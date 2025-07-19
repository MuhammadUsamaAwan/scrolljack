import { Link } from '@tanstack/react-router';
import { ScrollTextIcon } from 'lucide-react';
import { Button } from '~/components/ui/button';

export function Header() {
  return (
    <header className='border-b bg-background'>
      <div className='container mx-auto px-4'>
        <div className='flex h-16 items-center justify-between'>
          <div className='flex items-center space-x-8'>
            <Link to='/' className='flex items-center gap-1 font-bold text-xl'>
              <ScrollTextIcon className='h-6 w-6' />
              Scrolljack
            </Link>
            <nav className='flex space-x-4'>
              <Button variant='ghost' asChild>
                <Link to='/'>Home</Link>
              </Button>
              <Button variant='ghost' asChild>
                <Link to='/modlists'>Modlists</Link>
              </Button>
            </nav>
          </div>
        </div>
      </div>
    </header>
  );
}
