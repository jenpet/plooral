# Refactorings
- JSON representation to actual object vs. shape configuration (StickyNoteConfig). The UI should rather build on the JSON representation instead of wrapping it with custom logic to keep the compatibility with Konva.
- Add a separate layer filled with a single rect that takes all the space to provide a background color