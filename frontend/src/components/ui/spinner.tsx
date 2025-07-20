import { Loader2 } from 'lucide-react';
import { cn } from '~/lib/utils';

export function Spinner({ className }: { className?: string }) {
  return (
    <div className={cn('flex items-center justify-center', className)}>
      <Loader2 className='animate-spin h-6 w-6' />
    </div>
  );
}
