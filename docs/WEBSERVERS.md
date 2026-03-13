# Web Servers Guide

Lumine supports three popular web servers: Nginx, Apache, and Caddy.

## Supported Web Servers

### Nginx

**Description:** High-performance web server and reverse proxy

**Pros:**
- Very fast and lightweight
- Excellent for static files
- Great reverse proxy
- Low memory usage

**Cons:**
- Configuration can be complex
- No .htaccess support

**Default Port:** 80

**Configuration Example:**
```nginx
server {
    listen 80;
    server_name myapp.test;
    root /var/www/myapp/public;

    index index.html index.php;

    location / {
        try_files $uri $uri/ /index.php?$query_string;
    }

    location ~ \.php$ {
        fastcgi_pass php:9000;
        fastcgi_index index.php;
        fastcgi_param SCRIPT_FILENAME $document_root$fastcgi_script_name;
        include fastcgi_params;
    }
}
```

### Apache

**Description:** Most popular web server with extensive module support

**Pros:**
- .htaccess support
- Extensive documentation
- Many modules available
- Familiar to most developers

**Cons:**
- Higher memory usage
- Slower than Nginx for static files

**Default Port:** 80

**Configuration Example:**
```apache
<VirtualHost *:80>
    ServerName myapp.test
    DocumentRoot /var/www/myapp/public

    <Directory /var/www/myapp/public>
        Options Indexes FollowSymLinks
        AllowOverride All
        Require all granted
    </Directory>

    <FilesMatch \.php$>
        SetHandler "proxy:fcgi://php:9000"
    </FilesMatch>
</VirtualHost>
```

### Caddy

**Description:** Modern web server with automatic HTTPS

**Pros:**
- Automatic HTTPS
- Simple configuration
- Built-in reverse proxy
- Modern and actively developed

**Cons:**
- Newer, less documentation
- Different configuration syntax

**Default Ports:** 8085 (HTTP), 8445 (HTTPS)

**Configuration Example:**
```caddy
myapp.test {
    root * /var/www/myapp/public
    php_fastcgi php:9000
    file_server
    encode gzip
}
```

## Setup

### Start Web Server

```bash
# Via Makefile
make db-setup

# Or manually
docker compose -f docker-compose.db.yml up -d nginx
docker compose -f docker-compose.db.yml up -d apache
docker compose -f docker-compose.db.yml up -d caddy
```

### Access

- Nginx: http://localhost:80
- Apache: http://localhost:80
- Caddy: http://localhost:8085

## Configuration

### Nginx Configuration

Location: `~/.lumine/nginx/`

```bash
# Create site config
cat > ~/.lumine/nginx/myapp.conf << EOF
server {
    listen 80;
    server_name myapp.test;
    root /var/www/myapp/public;
    index index.php index.html;

    location / {
        try_files \$uri \$uri/ /index.php?\$query_string;
    }

    location ~ \.php$ {
        fastcgi_pass php:9000;
        fastcgi_index index.php;
        fastcgi_param SCRIPT_FILENAME \$document_root\$fastcgi_script_name;
        include fastcgi_params;
    }
}
EOF

# Reload Nginx
docker exec lumine-nginx nginx -s reload
```

### Apache Configuration

Location: `~/.lumine/apache/`

```bash
# Create site config
cat > ~/.lumine/apache/myapp.conf << EOF
<VirtualHost *:80>
    ServerName myapp.test
    DocumentRoot /var/www/myapp/public

    <Directory /var/www/myapp/public>
        AllowOverride All
        Require all granted
    </Directory>
</VirtualHost>
EOF

# Reload Apache
docker exec lumine-apache apachectl graceful
```

### Caddy Configuration

Location: `scripts/caddy/Caddyfile`

```bash
# Edit Caddyfile
cat >> scripts/caddy/Caddyfile << EOF
myapp.test {
    root * /var/www/myapp/public
    php_fastcgi php:9000
    file_server
    encode gzip
}
EOF

# Reload Caddy
docker exec lumine-caddy caddy reload --config /etc/caddy/Caddyfile
```

## Use Cases

### Static Site

**Nginx:**
```nginx
server {
    listen 80;
    server_name static.test;
    root /var/www/static;
    index index.html;
}
```

**Apache:**
```apache
<VirtualHost *:80>
    ServerName static.test
    DocumentRoot /var/www/static
</VirtualHost>
```

**Caddy:**
```caddy
static.test {
    root * /var/www/static
    file_server
}
```

### PHP Application (Laravel)

**Nginx:**
```nginx
server {
    listen 80;
    server_name laravel.test;
    root /var/www/laravel/public;

    location / {
        try_files $uri $uri/ /index.php?$query_string;
    }

    location ~ \.php$ {
        fastcgi_pass php:9000;
        fastcgi_param SCRIPT_FILENAME $document_root$fastcgi_script_name;
        include fastcgi_params;
    }
}
```

**Apache:**
```apache
<VirtualHost *:80>
    ServerName laravel.test
    DocumentRoot /var/www/laravel/public

    <Directory /var/www/laravel/public>
        AllowOverride All
        Require all granted
    </Directory>

    <FilesMatch \.php$>
        SetHandler "proxy:fcgi://php:9000"
    </FilesMatch>
</VirtualHost>
```

**Caddy:**
```caddy
laravel.test {
    root * /var/www/laravel/public
    php_fastcgi php:9000
    file_server
}
```

### Reverse Proxy (Next.js)

**Nginx:**
```nginx
server {
    listen 80;
    server_name nextjs.test;

    location / {
        proxy_pass http://localhost:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
    }
}
```

**Apache:**
```apache
<VirtualHost *:80>
    ServerName nextjs.test

    ProxyPreserveHost On
    ProxyPass / http://localhost:3000/
    ProxyPassReverse / http://localhost:3000/
</VirtualHost>
```

**Caddy:**
```caddy
nextjs.test {
    reverse_proxy localhost:3000
}
```

## Comparison

| Feature | Nginx | Apache | Caddy |
|---------|-------|--------|-------|
| Performance | Excellent | Good | Excellent |
| Memory Usage | Low | Medium | Low |
| Configuration | Complex | Medium | Simple |
| .htaccess | No | Yes | No |
| Auto HTTPS | No | No | Yes |
| Reverse Proxy | Excellent | Good | Excellent |
| Static Files | Excellent | Good | Excellent |
| PHP Support | Via FastCGI | Native/FastCGI | Via FastCGI |
| Learning Curve | Steep | Medium | Easy |

## Best Practices

### 1. Choose Based on Needs

- **Nginx**: High traffic, reverse proxy, static files
- **Apache**: .htaccess needed, familiar environment
- **Caddy**: Modern apps, automatic HTTPS, simple config

### 2. Use FastCGI for PHP

All three servers work best with PHP-FPM:

```yaml
services:
  - name: php
    type: php
    version: 8.2-fpm
    port: 9000
```

### 3. Enable Compression

**Nginx:**
```nginx
gzip on;
gzip_types text/plain text/css application/json;
```

**Apache:**
```apache
<IfModule mod_deflate.c>
    AddOutputFilterByType DEFLATE text/html text/plain text/css
</IfModule>
```

**Caddy:**
```caddy
encode gzip
```

### 4. Set Up Logging

**Nginx:**
```nginx
access_log /var/log/nginx/access.log;
error_log /var/log/nginx/error.log;
```

**Apache:**
```apache
ErrorLog ${APACHE_LOG_DIR}/error.log
CustomLog ${APACHE_LOG_DIR}/access.log combined
```

**Caddy:**
```caddy
log {
    output file /var/log/caddy/access.log
}
```

## Troubleshooting

### Nginx

```bash
# Test configuration
docker exec lumine-nginx nginx -t

# Reload
docker exec lumine-nginx nginx -s reload

# View logs
docker logs lumine-nginx
```

### Apache

```bash
# Test configuration
docker exec lumine-apache apachectl configtest

# Reload
docker exec lumine-apache apachectl graceful

# View logs
docker logs lumine-apache
```

### Caddy

```bash
# Validate configuration
docker exec lumine-caddy caddy validate --config /etc/caddy/Caddyfile

# Reload
docker exec lumine-caddy caddy reload --config /etc/caddy/Caddyfile

# View logs
docker logs lumine-caddy
```

## Migration

### From Apache to Nginx

1. Convert .htaccess rules to Nginx location blocks
2. Update PHP handler from mod_php to FastCGI
3. Test configuration
4. Switch DNS/proxy

### From Nginx to Caddy

1. Convert server blocks to Caddy sites
2. Simplify configuration (Caddy handles many things automatically)
3. Test configuration
4. Switch

## See Also

- [Configuration Guide](CONFIGURATION.md)
- [Project Setup](QUICKSTART.md)
- [Domain Management](PORT_MANAGEMENT.md)
