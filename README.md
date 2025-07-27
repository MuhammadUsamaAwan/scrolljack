# 📜 Scrolljack

Scrolljack is a desktop application built with [Wails](https://wails.io/) that transforms your <code>.wabbajack</code> files into readable modding guides

## Installation

- **For Windows:** [Download EXE](https://github.com/MuhammadUsamaAwan/scrolljack/releases/download/v0.1.0/scrolljack.exe)
- **For Linux:** [Download Binary](https://github.com/MuhammadUsamaAwan/scrolljack/releases/download/v0.1.0/scrolljack)

## Installation (Development)

### Prerequisites

- [Go 1.23+](https://go.dev/)
- CGO enabled
- [Node.js 22+](https://nodejs.org/)
- [Wails CLI](https://wails.io/)

### Local Development

```bash
git clone https://github.com/MuhammadUsamaAwan/scrolljack.git
cd scrolljack
cd frontend
pnpm install
cd ..
wails dev
```

## How to Use

1. Launch Scrolljack.
2. Select your <code>.wabbajack</code> file.
3. The import process may take a few seconds to several minutes depending on the size of the modlist and your CPU.
4. Once you see the **Modlist import completed** message, navigate to the **Modlists** page.
5. On this page, you can; View your imported modlists, Search for modlists, Delete modlists.
6. Click on a modlist to open its **Details Page**; Switch between profiles, Download profile files, Browse mods organized by separators.
7. Clicking on the mod will reveal it's archive(s) with links and **Show/Hide Files** files button.
8. You can download individual **Inline** and **RemappedInline** files.
9. If a file is marked as **PatchedFromArchive**, you can apply its patch; Select the original file, The patch is applied and saved to your Downloads folder, For text files, you'll see a diff view, For plugin or texture changes, use tools like xEdit or NIFviewer.
10. Use the **Detect FOMOD Options** button under each mod: Wabbajack doesn't expose which mods have FOMODs, Select the archive to scan, Detection may take time depending on file size, Results show a list of possible options with confidence scores.
11. Data Location: On Windows; <code>%APPDATA%/Roaming/scrolljack</code> On Linux: <code>~/.config/scrolljack</code>
