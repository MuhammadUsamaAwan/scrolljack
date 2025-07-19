import { createFileRoute } from '@tanstack/react-router';
import { useEffect, useState } from 'react';
import { toast } from 'sonner';
import { Hero } from '~/components/hero';
import { Button } from '~/components/ui/button';
import { ProcessWabbajackFile } from '~/wailsjs/go/main/App';
import { EventsOn } from "~/wailsjs/runtime";


export const Route = createFileRoute('/')({
  component: RouteComponent,
});

function RouteComponent() {
   const [progress, setProgress] = useState<string[]>([]);

    useEffect(() => {
        EventsOn("progress_update", (data) => {
            setProgress(prev => [...prev, data]);
        });
    }, []);



  return (
    <main className='container mx-auto space-y-8 px-4 py-10'>
      <Hero />
      <div className='flex justify-center'>
        <Button
          size='lg'
          onClick={async () => {
            setProgress([]);
            try {
              await ProcessWabbajackFile();
            } catch (error) {
              toast.error(`Error processing Wabbajack file: ${error instanceof Error ? error.message : 'Unknown error'}`);
            }
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
    </main>
  );
}
