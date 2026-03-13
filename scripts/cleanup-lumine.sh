#!/bin/bash
# Lumine Cleanup Script
# This script safely cleans up Lumine containers and optionally data

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

echo -e "${CYAN}🌟 Lumine Cleanup Script${NC}"
echo ""

# Check if Docker is running
if ! docker ps > /dev/null 2>&1; then
    echo -e "${RED}❌ Docker is not running!${NC}"
    exit 1
fi

# Menu
echo "Select cleanup option:"
echo "  1) Stop containers only"
echo "  2) Remove containers (keep data)"
echo "  3) Remove containers + volumes (DELETE DATA)"
echo "  4) Nuclear cleanup (REMOVE EVERYTHING)"
echo "  5) Cancel"
echo ""
read -p "Enter choice [1-5]: " choice

case $choice in
    1)
        echo -e "${CYAN}Stopping containers...${NC}"
        docker ps -q --filter "name=lumine-" | xargs -r docker stop
        echo -e "${GREEN}✓ Containers stopped!${NC}"
        ;;
    2)
        echo -e "${YELLOW}⚠️  This will remove all Lumine containers${NC}"
        read -p "Continue? [y/N] " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            echo -e "${CYAN}Stopping containers...${NC}"
            docker ps -q --filter "name=lumine-" | xargs -r docker stop
            echo -e "${CYAN}Removing containers...${NC}"
            docker ps -aq --filter "name=lumine-" | xargs -r docker rm
            echo -e "${GREEN}✓ Containers removed!${NC}"
        fi
        ;;
    3)
        echo -e "${RED}⚠️  WARNING: This will DELETE ALL DATABASE DATA!${NC}"
        read -p "Type 'yes' to confirm: " confirm
        if [ "$confirm" = "yes" ]; then
            # Backup first
            echo -e "${CYAN}Creating backup...${NC}"
            BACKUP_FILE="lumine-backup-$(date +%Y%m%d-%H%M%S).sql"
            docker exec lumine-mysql mysqldump -u root -proot --all-databases > "$BACKUP_FILE" 2>/dev/null || true
            if [ -f "$BACKUP_FILE" ]; then
                echo -e "${GREEN}✓ Backup saved: $BACKUP_FILE${NC}"
            fi
            
            echo -e "${CYAN}Stopping containers...${NC}"
            docker ps -q --filter "name=lumine-" | xargs -r docker stop
            echo -e "${CYAN}Removing containers...${NC}"
            docker ps -aq --filter "name=lumine-" | xargs -r docker rm
            echo -e "${CYAN}Removing volumes...${NC}"
            docker volume ls -q --filter "name=lumine_" | xargs -r docker volume rm
            echo -e "${GREEN}✓ Complete cleanup done!${NC}"
        else
            echo -e "${YELLOW}Cleanup cancelled.${NC}"
        fi
        ;;
    4)
        echo -e "${RED}☢️  NUCLEAR OPTION: This will DESTROY EVERYTHING!${NC}"
        echo -e "${RED}   - All containers${NC}"
        echo -e "${RED}   - All volumes (data)${NC}"
        echo -e "${RED}   - Network${NC}"
        echo -e "${RED}   - Docker cache${NC}"
        echo ""
        read -p "Type 'DESTROY' to confirm: " confirm
        if [ "$confirm" = "DESTROY" ]; then
            # Backup first
            echo -e "${CYAN}Creating backup...${NC}"
            BACKUP_FILE="lumine-backup-$(date +%Y%m%d-%H%M%S).sql"
            docker exec lumine-mysql mysqldump -u root -proot --all-databases > "$BACKUP_FILE" 2>/dev/null || true
            if [ -f "$BACKUP_FILE" ]; then
                echo -e "${GREEN}✓ Backup saved: $BACKUP_FILE${NC}"
            fi
            
            echo -e "${CYAN}Stopping containers...${NC}"
            docker ps -q --filter "name=lumine-" | xargs -r docker stop
            echo -e "${CYAN}Removing containers...${NC}"
            docker ps -aq --filter "name=lumine-" | xargs -r docker rm
            echo -e "${CYAN}Removing volumes...${NC}"
            docker volume ls -q --filter "name=lumine_" | xargs -r docker volume rm
            echo -e "${CYAN}Removing network...${NC}"
            docker network rm lumine 2>/dev/null || true
            echo -e "${CYAN}Pruning Docker system...${NC}"
            docker system prune -af
            echo -e "${GREEN}✓ Nuclear cleanup complete!${NC}"
        else
            echo -e "${YELLOW}Destruction cancelled.${NC}"
        fi
        ;;
    5)
        echo -e "${YELLOW}Cleanup cancelled.${NC}"
        exit 0
        ;;
    *)
        echo -e "${RED}Invalid choice!${NC}"
        exit 1
        ;;
esac

echo ""
echo -e "${CYAN}Cleanup Summary:${NC}"
echo "Containers: $(docker ps -aq --filter 'name=lumine-' | wc -l)"
echo "Volumes: $(docker volume ls -q --filter 'name=lumine_' | wc -l)"
echo ""
echo -e "${GREEN}Done!${NC}"
