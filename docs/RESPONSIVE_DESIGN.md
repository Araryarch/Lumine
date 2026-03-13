# Lumine TUI - Responsive Design

## Overview
Lumine TUI sekarang fully responsive dan menyesuaikan dengan ukuran terminal.

## Height Calculation

### Terminal Height Breakdown
```
Total Height = Terminal Height
├─ Title Bar: 1 line
├─ Content Area: Height - 4 lines
│  ├─ Sidebar: Fixed height (matches content)
│  ├─ Main Panel: Fixed height (matches content)
│  └─ Detail Panel: Fixed height (matches content)
├─ Status Bar: 1 line
└─ Help Bar: 1 line
```

### Content Height Formula
```go
contentHeight = m.height - 4
// Where:
// - 1 line for title bar
// - 1 line for status bar
// - 1 line for help bar
// - 1 line for spacing
```

### Panel Height Formula
```go
panelHeight = m.height - 5
// This accounts for:
// - Title bar (1)
// - Status bar (1)
// - Help bar (1)
// - Panel padding (2)
```

## Width Calculation

### Terminal Width Breakdown
```
Total Width = Terminal Width
├─ Sidebar: 22 chars (fixed)
├─ Main Panel: (Width - 24) / 2
├─ Detail Panel: (Width - 24) / 2
└─ Margins: 2 chars total
```

### Panel Width Formula
```go
sidebarWidth = 22 // Fixed
mainPanelWidth = ((m.width - 24) / 2) - 2
detailPanelWidth = ((m.width - 24) / 2) - 2
```

## Fixed Elements

### Sidebar (Always Fixed)
- **Width**: 22 characters (never changes)
- **Height**: Matches content area height
- **Behavior**: Never scrolls, always visible
- **Content**: Menu items with icons
- **Max Items**: 10 items (fits in most terminals)

### Title Bar (Always Fixed)
- **Height**: 1 line
- **Width**: Full terminal width
- **Content**: Logo + Description
- **Style**: Blue background with white text

### Status Bar (Always Fixed)
- **Height**: 1 line
- **Width**: Full terminal width
- **Content**: View name, selection count, running status, messages
- **Style**: Surface color background

### Help Bar (Always Fixed)
- **Height**: 1 line
- **Width**: Full terminal width
- **Content**: Context-aware keybindings
- **Style**: Surface color background

## Scrollable Elements

### Main Panel Content
- **Vertical Scroll**: When content exceeds panel height
- **Horizontal**: Truncate with "..." if too long
- **Examples**:
  - Services list (when > 20 services)
  - Project types (when > 15 types)
  - Logs (scrollable with ↑/↓)

### Logs Panel
- **Scroll Behavior**: 
  - `↑` or `k`: Scroll up
  - `↓` or `j`: Scroll down
  - Auto-scroll to bottom on new logs
- **Visible Lines**: `panelHeight - 6`
- **Scroll Indicator**: Shows "↑ X more logs above" / "↓ X more logs below"

## Minimum Requirements

### Minimum Terminal Size
```
Width:  80 characters (recommended: 120+)
Height: 24 lines (recommended: 30+)
```

### Optimal Terminal Size
```
Width:  120-140 characters
Height: 30-40 lines
```

### Behavior at Minimum Size
- Sidebar: Still 22 chars (fixed)
- Main Panel: ~27 chars (may truncate content)
- Detail Panel: ~27 chars (may truncate content)
- All panels maintain structure

## Responsive Behavior

### When Terminal Resizes

#### Width Changes
1. **Increase Width**:
   - Main and Detail panels expand proportionally
   - More content visible per line
   - Less truncation

2. **Decrease Width**:
   - Main and Detail panels shrink proportionally
   - More truncation with "..."
   - Sidebar remains 22 chars

#### Height Changes
1. **Increase Height**:
   - All panels expand vertically
   - More items visible without scrolling
   - Better spacing

2. **Decrease Height**:
   - All panels shrink vertically
   - More scrolling required
   - Content remains accessible

### Content Adaptation

#### Long Text Handling
```go
// Truncate if too long
if lipgloss.Width(text) > maxWidth {
    text = text[:maxWidth-3] + "..."
}
```

#### List Items
- Show as many as fit in panel height
- Scroll indicators when more items exist
- Cursor navigation works regardless of visible items

#### Empty States
- Centered in panel
- Helpful messages
- Action hints

## Panel Spacing

### Internal Padding
```
All Panels:
├─ Top: 1 line
├─ Right: 2 chars
├─ Bottom: 1 line
└─ Left: 2 chars
```

### Between Panels
```
Sidebar → Main: 0 chars (borders touch)
Main → Detail: 0 chars (borders touch)
```

### Content Spacing
```
Header: 2 lines after (includes underline)
Sections: 1 line between
Items: 0 lines between (consecutive)
Footer: 2 lines before
```

## Border Behavior

### Active Panel
- **Border**: ThickBorder (`╔═╗║║╚═╝`)
- **Color**: Primary Blue (#89b4fa)
- **Width**: 2 chars per side

### Inactive Panel
- **Border**: RoundedBorder (`╭─╮││╰─╯`)
- **Color**: Surface1 (#45475a)
- **Width**: 1 char per side

### Border Calculation
```go
// Content width includes borders
contentWidth = panelWidth - 4 // 2 chars per side
contentHeight = panelHeight - 2 // 1 line top + bottom
```

## Overflow Handling

### Vertical Overflow
1. **Lists**: Scrolling with indicators
2. **Text**: Wrap or truncate based on context
3. **Logs**: Scrollable with offset tracking

### Horizontal Overflow
1. **Text**: Truncate with "..."
2. **Badges**: Fixed width, never truncate
3. **Icons**: Fixed width (1-2 chars)

## Testing Different Sizes

### Small Terminal (80x24)
```
✓ All panels visible
✓ Sidebar fixed at 22 chars
✓ Main/Detail ~27 chars each
✓ Content may truncate
✓ Scrolling works
```

### Medium Terminal (100x30)
```
✓ All panels comfortable
✓ Sidebar fixed at 22 chars
✓ Main/Detail ~37 chars each
✓ Less truncation
✓ Better spacing
```

### Large Terminal (120x40)
```
✓ All panels spacious
✓ Sidebar fixed at 22 chars
✓ Main/Detail ~47 chars each
✓ Minimal truncation
✓ Optimal experience
```

### Extra Large Terminal (140x50)
```
✓ All panels very spacious
✓ Sidebar fixed at 22 chars
✓ Main/Detail ~57 chars each
✓ No truncation
✓ Excellent spacing
```

## Implementation Details

### Height Tracking
```go
type model struct {
    width  int  // Updated on tea.WindowSizeMsg
    height int  // Updated on tea.WindowSizeMsg
    // ...
}
```

### Update on Resize
```go
case tea.WindowSizeMsg:
    m.width = msg.Width
    m.height = msg.Height
    m.ready = true
    return m, nil
```

### Panel Rendering
```go
// All panels use consistent height
panelHeight := m.height - 5

// Sidebar uses content height
contentHeight := m.height - 4
sidebar := m.renderSidebarFixed(22, contentHeight)
```

## Best Practices

### For Developers
1. Always use `m.height` and `m.width` for calculations
2. Never hardcode panel dimensions
3. Test with different terminal sizes
4. Handle overflow gracefully
5. Provide scroll indicators

### For Users
1. Use terminal size 120x30 or larger for best experience
2. Resize terminal anytime - TUI adapts automatically
3. Use scrolling (↑/↓) for long lists
4. Check help bar for context-specific keys

## Known Limitations

1. **Minimum Width**: Below 80 chars, layout may break
2. **Minimum Height**: Below 24 lines, content may overlap
3. **Very Small Terminals**: Some content will be truncated
4. **No Horizontal Scroll**: Content truncates instead

## Future Improvements

1. **Dynamic Sidebar**: Collapse to icons only on small screens
2. **Horizontal Scroll**: For very long lines
3. **Zoom Levels**: Adjust font size/spacing
4. **Layout Modes**: Single panel mode for small terminals
5. **Responsive Badges**: Shorter versions for small screens
