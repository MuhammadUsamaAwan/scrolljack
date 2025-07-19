import { createFileRoute } from '@tanstack/react-router';
import { Hero } from '~/components/hero';
import { Button } from '~/components/ui/button';
import { SelectFile } from '~/wailsjs/go/main/App';

export const Route = createFileRoute('/')({
  component: RouteComponent,
});

function RouteComponent() {
  return (
    <main className='container mx-auto space-y-8 px-4 py-10'>
      <Hero />
      <div className='flex justify-center'>
        <Button
          size='lg'
          onClick={() => {
            SelectFile();
          }}
        >
          Select a Wabbajack file
        </Button>
      </div>
    </main>
  );
}
