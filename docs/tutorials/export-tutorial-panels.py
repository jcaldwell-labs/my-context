#!/usr/bin/env python3
"""
Export Tutorial Panels Using Wonderings TUI

Uses the Wonderings project's TUI modules to export context explorer and detail panels
as dark-mode styled HTML files for each tutorial.
"""

import sys
import os
import re
from pathlib import Path
from datetime import datetime

# Add Wonderings project to path to import TUI modules
WONDERINGS_PATH = Path("/home/be-dev-agent/projects/wonderings")
sys.path.insert(0, str(WONDERINGS_PATH))

# Import Wonderings TUI modules
try:
    from tui_modules import get_all_contexts, get_active_context, get_all_projects
    from tui_modules.explorer_panel import create_explorer_panel
    from tui_modules.detail_panel import create_detail_panel
    from rich.console import Console
except ImportError as e:
    print(f"❌ Error: Could not import Wonderings TUI modules: {e}")
    print(f"Make sure Wonderings project exists at: {WONDERINGS_PATH}")
    sys.exit(1)

# Directories
CONTEXT_HOMES_DIR = Path(__file__).parent / "context-homes"
TUTORIALS_DIR = Path(__file__).parent

# Tutorial context home mapping
TUTORIALS = {
    "01": {
        "name": "Backend Developer Solo",
        "context_homes": [
            ("tutorial-01-backend-solo", "Alice - Payment Retry Logic")
        ]
    },
    "02": {
        "name": "Frontend Developer Solo",
        "context_homes": [
            ("tutorial-02-frontend-solo", "Bob - Checkout UI")
        ]
    },
    "03": {
        "name": "QA Engineer Solo",
        "context_homes": [
            ("tutorial-03-qa-solo", "Carol - Payment Testing")
        ]
    },
    "04": {
        "name": "Multi-Project Consultant",
        "context_homes": [
            ("tutorial-04-multi-project", "Alice - 3 Client Projects")
        ]
    },
    "05": {
        "name": "Scrum Master Sprint Management",
        "context_homes": [
            ("tutorial-05-scrum-master", "Dave - Sprint 5")
        ]
    },
    "06": {
        "name": "Team Handoff",
        "context_homes": [
            ("tutorial-06-team-alice", "Alice - Backend API"),
            ("tutorial-06-team-bob", "Bob - Frontend Integration")
        ]
    },
    "07": {
        "name": "Signal Coordination",
        "context_homes": [
            ("tutorial-07-release-alice", "Alice - Backend Release"),
            ("tutorial-07-release-bob", "Bob - Frontend Release"),
            ("tutorial-07-release-carol", "Carol - QA Testing"),
            ("tutorial-07-release-eve", "Eve - Product Coordination")
        ]
    },
    "08": {
        "name": "Agent Workflows",
        "context_homes": [
            ("tutorial-08-human-alice", "Alice - OAuth Feature"),
            ("tutorial-08-agent-claude", "Claude Agent - Code Assistance"),
            ("tutorial-08-agent-cicd", "CI/CD Agent - Build & Test"),
            ("tutorial-08-agent-qa", "QA Bot - E2E Testing")
        ]
    }
}

def extract_rich_styles(html_content):
    """Extract style definitions from Rich HTML"""
    match = re.search(r'<style>(.*?)</style>', html_content, re.DOTALL)
    if match:
        styles = match.group(1)
        # Remove the body styling to avoid light background
        styles = re.sub(r'body\s*\{[^}]*\}', '', styles)
        return styles
    return ""

def extract_html_body(html_content):
    """Extract just the body content from Rich HTML"""
    match = re.search(r'<pre[^>]*>(.*?)</pre>', html_content, re.DOTALL)
    if match:
        return f"<pre>{match.group(1)}</pre>"
    return "<pre>Export data not found</pre>"

def create_dark_mode_html(title, rich_styles, body_content, context_home_name, metadata):
    """Create dark-mode HTML wrapper"""
    return f"""<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>{title}</title>
<style>
/* Dark mode styling */
body {{
    background-color: #1a1a1a;
    color: #e0e0e0;
    font-family: 'Menlo', 'DejaVu Sans Mono', 'Courier New', monospace;
    margin: 0;
    padding: 20px;
    line-height: 1.4;
}}

pre {{
    background-color: #0d0d0d;
    color: #e0e0e0;
    padding: 20px;
    border-left: 2px solid #444;
    border-radius: 4px 0 0 4px;
    overflow-x: auto;
    font-size: 13px;
}}

code {{
    font-family: inherit;
    color: inherit;
}}

/* Preserve all the color styles from Rich */
{rich_styles}

/* Make sure text is readable */
.r1, .r2, .r3, .r4, .r5, .r6, .r7, .r8, .r9, .r10, .r11, .r12,
.r13, .r14, .r15, .r16, .r17, .r18, .r19, .r20 {{
    font-weight: normal;
}}

/* Dim text should still be readable on dark background */
[style*="color: #7f7f7f"],
[style*="color: #808080"] {{
    color: #888888 !important;
}}

/* Ensure active items show properly */
[style*="background-color: #000080"] {{
    background-color: #1a3a5c !important;
}}

/* Links and interactive elements */
a {{
    color: #66d9ef;
    text-decoration: none;
}}

a:hover {{
    text-decoration: underline;
}}

/* Tutorial metadata */
.metadata {{
    background-color: #0d0d0d;
    padding: 15px;
    margin-bottom: 20px;
    border-left: 3px solid #66d9ef;
    border-radius: 4px;
}}

.metadata h3 {{
    margin-top: 0;
    color: #66d9ef;
}}

.metadata p {{
    margin: 5px 0;
    color: #e0e0e0;
}}
</style>
</head>
<body>
<div class="metadata">
<h3>{title}</h3>
<p><strong>Context Home:</strong> {context_home_name}</p>
<p><strong>Exported:</strong> {metadata['timestamp']}</p>
<p><strong>Total Contexts:</strong> {metadata['total_contexts']}</p>
<p><strong>Projects:</strong> {metadata['total_projects']}</p>
<p><strong>Active Context:</strong> {metadata['active_context']}</p>
</div>
{body_content}
</body>
</html>"""

def export_explorer_panel(context_home_dir, output_file):
    """Export explorer panel for a context home"""
    # Set environment for this context home
    os.environ["MY_CONTEXT_HOME"] = str(context_home_dir)

    # Get data
    contexts = get_all_contexts()
    active = get_active_context()
    projects = get_all_projects()

    # Create export console
    export_console = Console(record=True, width=140, force_terminal=True, legacy_windows=False)

    # Render explorer panel
    export_console.print("\n[bold cyan]CONTEXT HIERARCHY PANEL[/bold cyan]")
    export_console.print(f"[dim]Exported: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}[/dim]")
    export_console.print(f"[dim]Context Home: {context_home_dir.name}[/dim]")
    export_console.print(f"[dim]Total Contexts: {len(contexts)} | Projects: {len(projects)} | Active: {active.get('name') if active else 'none'}[/dim]\n")

    explorer = create_explorer_panel(page=0, show_all=True)
    export_console.print(explorer)

    # Get HTML
    html_content = export_console.export_html()

    # Extract styles and body
    rich_styles = extract_rich_styles(html_content)
    body_content = extract_html_body(html_content)

    # Metadata
    metadata = {
        'timestamp': datetime.now().strftime('%Y-%m-%d %H:%M:%S'),
        'total_contexts': len(contexts),
        'total_projects': len(projects),
        'active_context': active.get('name') if active else 'none'
    }

    # Create dark-mode HTML
    dark_mode_html = create_dark_mode_html(
        title=f"Explorer Panel - {context_home_dir.name}",
        rich_styles=rich_styles,
        body_content=body_content,
        context_home_name=context_home_dir.name,
        metadata=metadata
    )

    # Write to file
    output_file.write_text(dark_mode_html, encoding="utf-8")
    return len(contexts)

def export_detail_panel(context_home_dir, output_file):
    """Export detail panel for active context in a context home"""
    # Set environment for this context home
    os.environ["MY_CONTEXT_HOME"] = str(context_home_dir)

    # Get active context
    active = get_active_context()

    if not active:
        print(f"    ⚠️  No active context in {context_home_dir.name}, skipping detail export")
        return 0

    # Create export console
    export_console = Console(record=True, width=140, force_terminal=True, legacy_windows=False)

    # Render detail panel
    export_console.print("\n[bold cyan]CONTEXT DETAIL PANEL[/bold cyan]")
    export_console.print(f"[dim]Exported: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}[/dim]")
    export_console.print(f"[dim]Context: {active.get('name')}[/dim]\n")

    detail = create_detail_panel()
    export_console.print(detail)

    # Get HTML
    html_content = export_console.export_html()

    # Extract styles and body
    rich_styles = extract_rich_styles(html_content)
    body_content = extract_html_body(html_content)

    # Metadata
    metadata = {
        'timestamp': datetime.now().strftime('%Y-%m-%d %H:%M:%S'),
        'total_contexts': 1,
        'total_projects': 0,
        'active_context': active.get('name')
    }

    # Create dark-mode HTML
    dark_mode_html = create_dark_mode_html(
        title=f"Detail Panel - {active.get('name')}",
        rich_styles=rich_styles,
        body_content=body_content,
        context_home_name=context_home_dir.name,
        metadata=metadata
    )

    # Write to file
    output_file.write_text(dark_mode_html, encoding="utf-8")
    return 1

def export_all_tutorials():
    """Export panels for all tutorials"""
    print("=" * 70)
    print("MY-CONTEXT TUTORIAL PANEL EXPORT")
    print("=" * 70)
    print(f"Using Wonderings TUI from: {WONDERINGS_PATH}")
    print(f"Context homes: {CONTEXT_HOMES_DIR}")
    print()

    total_exported = 0

    for tutorial_num, tutorial_info in TUTORIALS.items():
        print(f"\n📘 Tutorial {tutorial_num}: {tutorial_info['name']}")

        tutorial_dir = TUTORIALS_DIR / f"tutorial-{tutorial_num}"
        tutorial_dir.mkdir(exist_ok=True)

        for context_home_name, description in tutorial_info['context_homes']:
            context_home_dir = CONTEXT_HOMES_DIR / context_home_name

            if not context_home_dir.exists():
                print(f"  ❌ Context home not found: {context_home_name}")
                continue

            print(f"  Exporting: {description}")

            # Export explorer panel
            explorer_file = tutorial_dir / f"{context_home_name}_explorer.html"
            try:
                count = export_explorer_panel(context_home_dir, explorer_file)
                print(f"    ✅ Explorer: {count} contexts → {explorer_file.name}")
                total_exported += 1
            except Exception as e:
                print(f"    ❌ Explorer export failed: {e}")

            # Export detail panel (if active context exists)
            detail_file = tutorial_dir / f"{context_home_name}_detail.html"
            try:
                count = export_detail_panel(context_home_dir, detail_file)
                if count > 0:
                    print(f"    ✅ Detail: {count} context → {detail_file.name}")
                    total_exported += 1
            except Exception as e:
                print(f"    ❌ Detail export failed: {e}")

    print()
    print("=" * 70)
    print(f"✅ PANEL EXPORT COMPLETE: {total_exported} files generated")
    print("=" * 70)
    print()
    print("Next step: Run build-tutorial-html.py to create tutorial pages")

if __name__ == "__main__":
    export_all_tutorials()
