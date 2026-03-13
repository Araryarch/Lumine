# TUI Improvements

## Overview
Lumine TUI telah ditingkatkan untuk memberikan pengalaman yang lebih profesional, mirip dengan tools seperti lazygit dan neovim.

## Key Improvements

### 1. **Visual Design - Catppuccin Mocha**
- **Catppuccin Color Scheme**: Theme yang nyaman dan modern
- **Better Borders**: Active panel menggunakan ThickBorder untuk membedakan dengan jelas
- **Status Icons**: Menggunakan Unicode icons (●, ☐, ☑, ▶, ⚡, ✓, ✗, ⚠) untuk visual yang lebih baik
- **Improved Typography**: Bold, italic, dan underline digunakan dengan tepat

### 2. **Dedicated Pages**
- **Services Page**: Manage Docker services
- **Projects Page**: Manage development projects
- **Databases Page**: Database management
- **Runtimes Page**: Runtime versions
- **Logs Page**: Real-time logs dengan scrolling
- **Background Tasks Page**: Monitor running operations
- **Settings Page**: Configuration management

### 3. **Logs System**
- **Real-time Logging**: Semua operasi dicatat dengan timestamp
- **Log Levels**: Info, Success, Error, Warning dengan color coding
- **Scrollable**: Navigate dengan ↑/↓
- **Service Tagging**: Setiap log ditag dengan service name
- **Persistent**: Menyimpan 1000 log terakhir

### 4. **Background Tasks**
- **Task Monitoring**: Lihat operasi yang berjalan di background
- **Status Tracking**: Running, Completed, Failed dengan icons
- **Progress Messages**: Update real-time untuk setiap task
- **History**: Menyimpan history task yang sudah selesai

### 2. **Input Components**
- **TextInput**: Component input text yang proper dengan cursor
- **NumberInput**: Input khusus untuk angka dengan validasi
- **Features**:
  - Cursor navigation (←/→, Home/End)
  - Character insertion/deletion
  - Placeholder text
  - Focus state visual
  - Border highlighting saat focused

### 3. **Status Bar**
- **View Indicator**: Menampilkan view aktif (Services, Projects, dll)
- **Selection Counter**: Jumlah item yang dipilih
- **Running Status**: Status services yang running dengan color coding
- **Messages**: Info messages dengan icon
- **Full Width**: Menggunakan seluruh lebar terminal

### 4. **Help Panel**
- **Context-Aware**: Help berubah sesuai panel aktif
- **Clear Keybindings**: Key + description yang jelas
- **Visual Separators**: Menggunakan │ untuk memisahkan commands
- **Color Coded**: Keys di-highlight dengan warna berbeda

### 5. **Title Bar**
- **Split Layout**: Logo di kiri, description di kanan
- **Full Width**: Menggunakan seluruh lebar terminal
- **Icon**: ⚡ icon untuk branding

### 6. **Service List**
- **Better Status Icons**: ● untuk running/stopped dengan color
- **Checkbox**: ☐/☑ untuk selection
- **Cursor**: ▶ untuk item aktif
- **Badges**: Color-coded badges untuk service types
- **Info Display**: Port dan version dengan color coding
- **Empty State**: Message saat belum ada services

### 7. **Keyboard Shortcuts**
- `q` atau `ctrl+c`: Quit
- `↑/↓` atau `j/k`: Navigate
- `h/l`: Switch panels
- `space`: Select/deselect
- `enter` atau `s`: Start service
- `x`: Stop service
- `r`: Restart service
- `v`: Change version
- `c`: Cleanup dialog
- `tab`: Next panel

## Color Palette (Catppuccin Mocha)

```
Primary:   #89b4fa (Blue)
Secondary: #cba6f7 (Mauve)
Success:   #a6e3a1 (Green)
Error:     #f38ba8 (Red)
Warning:   #f9e2af (Yellow)
Info:      #94e2d5 (Teal)
Muted:     #6c7086 (Overlay0)
Border:    #45475a (Surface1)
BG:        #1e1e2e (Base)
FG:        #cdd6f4 (Text)
Surface:   #313244 (Surface0)
Peach:     #fab387 (Peach)
Pink:      #f5c2e7 (Pink)
Lavender:  #b4befe (Lavender)
Sapphire:  #74c7ec (Sapphire)
```

## Future Improvements

### Planned Features
1. **Search/Filter**: Fuzzy search untuk services
2. **Logs Viewer**: Real-time logs dengan scrolling
3. **Resource Monitor**: CPU/Memory usage per container
4. **Multi-select Actions**: Bulk operations
5. **Configuration Editor**: Edit config dari TUI
6. **Themes**: Multiple color schemes
7. **Mouse Support**: Click to select/navigate
8. **Split Panes**: Multiple views simultaneously

### Input Forms
- Project creation form dengan proper inputs
- Service configuration form
- Port selection dengan validation
- Environment variables editor

### Advanced Features
- Command palette (seperti VSCode)
- Quick actions menu
- Notifications/toasts
- Progress indicators
- Confirmation dialogs yang lebih baik
- Context menus

## Technical Details

### Components Structure
```
internal/ui/
├── input.go           # Input components (TextInput, NumberInput)
├── styles.go          # All styles and colors
├── model.go           # Main model and state
├── panels.go          # Panel rendering
├── cleanup_dialog.go  # Cleanup dialog
├── version_selector.go # Version selector
└── ...
```

### Style Guidelines
- Use semantic colors (success, error, warning, info)
- Consistent padding and margins
- Clear visual hierarchy
- Responsive to terminal size
- Accessible color contrast

## Usage Examples

### Starting the TUI
```bash
lumine
```

### Navigation
- Use arrow keys or vim keys (j/k) to navigate
- Press `h` to go to sidebar, `l` to go to main panel
- Press `tab` to cycle through panels

### Managing Services
1. Navigate to a service
2. Press `space` to select (can select multiple)
3. Press `s` or `enter` to start
4. Press `x` to stop
5. Press `r` to restart

### Cleanup
1. Press `c` to open cleanup dialog
2. Use arrow keys to select cleanup option
3. Press `enter` to confirm
4. Press `esc` to cancel

## Contributing

When adding new UI components:
1. Follow the existing color scheme
2. Use semantic styles from `styles.go`
3. Ensure keyboard navigation works
4. Add help text for new keybindings
5. Test with different terminal sizes
