import { createFileRoute } from '@tanstack/react-router';
import { useEffect, useRef, useState } from 'react';
import { Hero } from '~/components/hero';
import { Button } from '~/components/ui/button';
import { ProcessWabbajackFile } from '~/wailsjs/go/main/App';
import { EventsOn } from '~/wailsjs/runtime';

export const Route = createFileRoute('/')({
  component: RouteComponent,
});

function RouteComponent() {
  const [progress, setProgress] = useState<string[]>([]);
  const bottomRef = useRef<HTMLDivElement | null>(null);

  useEffect(() => {
    EventsOn('progress_update', data => {
      setProgress(prev => [...prev, data]);
    });
  }, []);

  useEffect(() => {
    if (bottomRef.current) {
      bottomRef.current.scrollIntoView({ behavior: 'smooth' });
    }
  }, [progress]);

  return (
    <main className='container mx-auto space-y-8 px-4 py-10'>
      <Hero />
      <div className='flex justify-center'>
        <Button
          size='lg'
          onClick={async () => {
            setProgress([]);
            await ProcessWabbajackFile();
          }}
        >
          Select a Wabbajack file
        </Button>
      </div>
      {progress.length > 0 && (
        <div className='space-y-2 rounded-xl bg-card p-4 text-muted-foreground'>
          {progress.map(m => (
            <div key={m}>{m}</div>
          ))}
        </div>
      )}
      <div ref={bottomRef} />
    </main>
  );
}
