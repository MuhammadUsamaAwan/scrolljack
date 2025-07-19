import { createFileRoute } from '@tanstack/react-router';
import { Hero } from '~/components/hero';

export const Route = createFileRoute('/')({
  component: RouteComponent,
});

function RouteComponent() {
  return (
    <main className='container mx-auto space-y-8 px-4 py-10'>
      <Hero />
    </main>
  );
}
