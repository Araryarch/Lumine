# Lumine TUI Layout Guide

## Layout Structure

```
┌─────────────────────────────────────────────────────────────────────────┐
│ ⚡ LUMINE                                    Docker Development Manager │
├──────────┬──────────────────────────────┬──────────────────────────────┤
│          │                              │                              │
│ SIDEBAR  │      MAIN CONTENT            │      DETAIL PANEL            │
│ (Fixed)  │      (Dynamic)               │      (Context-aware)         │
│          │                              │                              │
│ ▶ ⚙ Services                            │                              │
│   📦 Projects                            │                              │
│   🗄 Databases                           │                              │
│   🔧 Runtimes                            │                              │
│   📋 Logs                                │                              │
│   ⚡ Tasks                                │                              │
│   ➕ New Project                          │                              │
│   ⚙ Settings                             │                              │
│   🔄 Refresh                              │                              │
│   🚪 Quit                                 │                              │
│          │                              │                              │
├──────────┴──────────────────────────────┴──────────────────────────────┤
│ ▶ Services  ✓ 2 selected  ● 3/5 running  ℹ Starting nginx...           │
├─────────────────────────────────────────────────────────────────────────┤
│ ↑/↓ navigate  space select  enter start  h/l panels  q quit            │
└─────────────────────────────────────────────────────────────────────────┘
```

## Panel Descriptions

### 1. Sidebar (Left - Fixed, 22 chars wide)
- **Always visible** - tidak scroll
- **Navigation menu** dengan icons
- **Active indicator** - `▶` untuk item yang dipilih
- **Highlight** - background color untuk selected item
- **Border** - ThickBorder saat active, RoundedBorder saat inactive

**Features:**
- Fixed width (22 characters)
- Fixed height (matches content area)
- Icons untuk setiap menu item
- Color-coded selection
- Keyboard navigation (↑/↓)

### 2. Main Content (Center - Dynamic)
- **Context-dependent** - berubah sesuai sidebar selection
- **Scrollable** - untuk content yang panjang
- **Interactive** - support selection, navigation
- **Border** - ThickBorder saat active

**Views:**
- **Services**: List services dengan status, port, version
- **Projects**: List projects dengan type, path
- **Databases**: Database connections
- **Runtimes**: Runtime versions (PHP, Node, Python, Rust)
- **Logs**: Real-time logs dengan scrolling
- **Tasks**: Background tasks dengan status
- **New Project**: Project creation wizard
- **Settings**: Configuration dan statistics

### 3. Detail Panel (Right - Context-aware)
- **Dynamic content** - berubah sesuai main content
- **Informative** - menampilkan detail dari selected item
- **Border** - ThickBorder saat active

**Content per View:**

#### Services View
```
╭─────────────────────────────╮
│ 📋 Service Details          │
│                             │
│ ╭─────────────────────────╮ │
│ │ Name: mysql             │ │
│ │ Type: MYSQL             │ │
│ │ Version: 8.0            │ │
│ │ Port: :3306             │ │
│ │ Status: ● Running       │ │
│ │ Container: a1b2c3d4e5f6 │ │
│ ╰─────────────────────────╯ │
│                             │
│ Environment Variables:      │
│ ╭─────────────────────────╮ │
│ │ MYSQL_ROOT_PASSWORD=*** │ │
│ │ MYSQL_DATABASE=app      │ │
│ ╰─────────────────────────╯ │
╰─────────────────────────────╯
```

#### Logs View
```
╭─────────────────────────────╮
│ 📊 Log Statistics           │
│                             │
│ ╭─────────────────────────╮ │
│ │ Total Logs: 127         │ │
│ │                         │ │
│ │ ✓ Success: 89           │ │
│ │ ✗ Errors: 12            │ │
│ │ ⚠ Warnings: 18          │ │
│ │ ℹ Info: 8               │ │
│ ╰─────────────────────────╯ │
╰─────────────────────────────╯
```

#### Tasks View
```
╭─────────────────────────────╮
│ 📊 Task Statistics          │
│                             │
│ ╭─────────────────────────╮ │
│ │ Total Tasks: 15         │ │
│ │                         │ │
│ │ ⟳ Running: 2            │ │
│ │ ✓ Completed: 11         │ │
│ │ ✗ Failed: 2             │ │
│ ╰─────────────────────────╯ │
╰─────────────────────────────╯
```

#### Settings View
```
╭─────────────────────────────╮
│ ℹ About                     │
│                             │
│ ╭─────────────────────────╮ │
│ │ Lumine                  │ │
│ │ Docker Development Mgr  │ │
│ │                         │ │
│ │ Version: 1.0.0          │ │
│ │ Theme: Catppuccin Mocha │ │
│ ╰─────────────────────────╯ │
╰─────────────────────────────╯
```

## Status Bar (Bottom)
- **View indicator**: Current view name dengan icon
- **Selection counter**: Jumlah item yang dipilih
- **Running status**: Services yang running
- **Messages**: Info/error messages
- **Full width**: Menggunakan seluruh lebar terminal

**Format:**
```
▶ Services  ✓ 2 selected  ● 3/5 running  ℹ Starting nginx...
```

## Help Bar (Bottom)
- **Context-aware**: Berubah sesuai active panel dan view
- **Keybindings**: Key + description
- **Visual separators**: `│` untuk memisahkan commands
- **Color-coded**: Keys di-highlight

**Examples:**

Services view:
```
↑/↓ navigate  space select  enter start  s start  x stop  r restart  v version  c cleanup  h/l panels  q quit
```

Sidebar:
```
↑/↓ navigate  enter select  l main panel  q quit
```

Logs view:
```
↑/↓ scroll  h/l panels  q quit
```

## Color Scheme (Catppuccin Mocha)

### Primary Colors
- **Blue** (#89b4fa): Primary actions, headers, active elements
- **Mauve** (#cba6f7): Secondary elements, icons
- **Green** (#a6e3a1): Success states, running services
- **Red** (#f38ba8): Errors, stopped services
- **Yellow** (#f9e2af): Warnings, ports
- **Teal** (#94e2d5): Info, versions

### Background Colors
- **Base** (#1e1e2e): Main background
- **Surface0** (#313244): Highlighted items
- **Surface1** (#45475a): Borders (inactive)

### Text Colors
- **Text** (#cdd6f4): Primary text
- **Overlay0** (#6c7086): Muted text, descriptions

## Border Styles

### Active Panel
```
╔═══════════════════════════╗
║                           ║
║   Active Panel Content    ║
║                           ║
╚═══════════════════════════╝
```
- **Style**: ThickBorder
- **Color**: Primary Blue (#89b4fa)

### Inactive Panel
```
╭───────────────────────────╮
│                           │
│  Inactive Panel Content   │
│                           │
╰───────────────────────────╯
```
- **Style**: RoundedBorder
- **Color**: Surface1 (#45475a)

## Navigation Flow

### Panel Navigation
1. **h** - Move to left panel (Sidebar)
2. **l** - Move to right panel (Main → Detail)
3. **tab** - Cycle through panels (Sidebar → Main → Detail → Sidebar)

### Sidebar Navigation
1. Select view dengan ↑/↓
2. Press Enter untuk switch ke view
3. Main panel akan update sesuai view
4. Detail panel akan update sesuai context

### Main Content Navigation
1. Navigate items dengan ↑/↓
2. Select/deselect dengan Space
3. Action dengan Enter atau shortcut keys
4. Detail panel akan update sesuai selected item

## Responsive Behavior

### Terminal Size
- **Minimum width**: 100 characters
- **Minimum height**: 24 lines
- **Optimal**: 120x30 or larger

### Panel Sizing
- **Sidebar**: Fixed 22 chars
- **Main Content**: (width - 24) / 2
- **Detail Panel**: (width - 24) / 2
- **Margins**: 2 chars between panels

### Content Overflow
- **Vertical**: Scrolling (logs, long lists)
- **Horizontal**: Truncation dengan "..."
- **Empty states**: Helpful messages

## Keyboard Shortcuts

### Global
- `q` or `ctrl+c`: Quit application
- `h`: Focus sidebar
- `l`: Focus main/detail panel
- `tab`: Next panel
- `esc`: Cancel/Go back

### Navigation
- `↑/↓` or `j/k`: Navigate items
- `enter`: Select/Confirm
- `space`: Toggle selection

### Services
- `s`: Start service
- `x`: Stop service
- `r`: Restart service
- `v`: Change version
- `c`: Cleanup dialog
- `delete`: Remove service
- `a`: Select all
- `d`: Deselect all

### Logs
- `↑/↓`: Scroll logs
- `g`: Go to top
- `G`: Go to bottom

## Best Practices

### Visual Hierarchy
1. **Headers**: Bold, underlined, colored
2. **Content**: Regular weight, readable
3. **Metadata**: Muted color, smaller emphasis
4. **Actions**: Highlighted, clear indicators

### Consistency
- Same border styles across panels
- Consistent spacing and padding
- Uniform icon usage
- Predictable keybindings

### Accessibility
- High contrast colors
- Clear visual indicators
- Keyboard-only navigation
- Screen reader friendly (text-based)

### Performance
- Efficient rendering
- Minimal redraws
- Smooth scrolling
- Responsive input
