import { diffLines } from 'diff';

export function FileDiff({ original, patched }: { original: string; patched: string }) {
  const diffs = diffLines(original, patched);

  return diffs.map((part, index) => (
    <pre
      key={index}
      className={`whitespace-pre-wrap ${
        part.added ? 'bg-green-100 text-green-700' : part.removed ? 'bg-red-100 text-red-700' : 'text-muted-foreground'
      }`}
    >
      {part.value}
    </pre>
  ));
}
