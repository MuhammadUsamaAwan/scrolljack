import { useQuery } from '@tanstack/react-query';
import { SearchIcon, XIcon } from 'lucide-react';
import { useMemo, useState } from 'react';
import { Mod } from '~/components/mod';
import { Button } from '~/components/ui/button';
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from '~/components/ui/collapsible';
import { Input } from '~/components/ui/input';
import { Spinner } from '~/components/ui/spinner';
import { profileModsQueryOptions } from '~/lib/query-options';

export function ProfileMods({ profileId }: { profileId: string }) {
  const { data, isPending } = useQuery(profileModsQueryOptions(profileId));
  const [searchTerm, setSearchTerm] = useState('');

  const filteredMods = useMemo(() => {
    if (!searchTerm.trim()) {
      return data;
    }
    const searchLower = searchTerm.toLowerCase();
    return data
      ?.map(group => ({
        ...group,
        mods: group.mods.filter(mod => mod.name.toLowerCase().includes(searchLower)),
      }))
      .filter(group => group.mods.length > 0);
  }, [data, searchTerm]);

  if (isPending) {
    return <Spinner />;
  }

  return (
    <div className='space-y-4'>
      <div className='relative'>
        <SearchIcon className='-translate-y-1/2 absolute top-1/2 left-3 h-4 w-4 text-muted-foreground' />
        <Input
          type='text'
          placeholder='Search mods...'
          value={searchTerm}
          onChange={e => setSearchTerm(e.target.value)}
          className='pl-10'
        />
        {searchTerm && (
          <Button
            variant='ghost'
            size='sm'
            onClick={() => setSearchTerm('')}
            className='-translate-y-1/2 absolute top-1/2 right-1 h-8 w-8 p-0'
            aria-label='Clear search'
          >
            <XIcon className='h-4 w-4' />
          </Button>
        )}
      </div>

      {filteredMods?.length === 0 && searchTerm ? (
        <div className='py-8 text-center text-muted-foreground'>No mods found matching "{searchTerm}"</div>
      ) : (
        filteredMods?.map(m => (
          <Collapsible key={m.separator}>
            <CollapsibleTrigger className='w-full cursor-pointer rounded-lg border bg-card px-4 py-2.5 font-semibold before:mr-1 before:inline-block before:text-muted-foreground before:text-xs before:duration-100 before:content-["⮞"] aria-expanded:before:rotate-90'>
              {m.separator}
              <span className='ml-2 text-muted-foreground text-sm'>
                ({m.mods.length} mod{m.mods.length !== 1 ? 's' : ''})
              </span>
            </CollapsibleTrigger>
            <CollapsibleContent className='mt-4 space-y-2'>
              {m.mods.map(mod => (
                <Mod key={mod.id} mod={mod} />
              ))}
            </CollapsibleContent>
          </Collapsible>
        ))
      )}
    </div>
  );
}
