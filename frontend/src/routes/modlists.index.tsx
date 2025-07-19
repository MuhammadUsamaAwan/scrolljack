import { useMutation, useSuspenseQuery } from '@tanstack/react-query';
import { createFileRoute, Link } from '@tanstack/react-router';
import { SearchIcon, Trash2, XIcon } from 'lucide-react';
import { useMemo, useState } from 'react';
import { toast } from 'sonner';
import { ModlistImage } from '~/components/modlist-image';
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from '~/components/ui/alert-dialog';
import { Badge } from '~/components/ui/badge';
import { Button } from '~/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '~/components/ui/card';
import { Input } from '~/components/ui/input';
import { queryClient } from '~/lib/query-client';
import { modListQueryOptions } from '~/lib/query-options';
import { pascalCaseToTitleCase } from '~/lib/utils';
import { DeleteModlist } from '~/wailsjs/go/main/App';

export const Route = createFileRoute('/modlists/')({
  loader: async () => {
    const modlists = queryClient.ensureQueryData(modListQueryOptions);
    return { modlists };
  },
  component: RouteComponent,
});

function RouteComponent() {
  const { data: initialModlists } = useSuspenseQuery(modListQueryOptions);
  const { mutateAsync } = useMutation({
    mutationFn: DeleteModlist,
    meta: {
      invalidateQueries: modListQueryOptions.queryKey,
    },
  });
  const [searchTerm, setSearchTerm] = useState('');

  const filteredModlists = useMemo(
    () =>
      (initialModlists ?? []).filter(
        modlist =>
          modlist.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
          modlist.author?.toLowerCase().includes(searchTerm.toLowerCase()) ||
          modlist.description?.toLowerCase().includes(searchTerm.toLowerCase()) ||
          modlist.game_type?.toLowerCase().includes(searchTerm.toLowerCase())
      ),
    [initialModlists, searchTerm]
  );

  const handleDeleteModlist = async (modlistId: string) => {
    await mutateAsync(modlistId);
    toast.promise(mutateAsync(modlistId), {
      loading: 'Deleting modlist this may take a while...',
      success: () => 'Modlist deleted successfully',
      error: error => `Failed to delete modlist: ${error instanceof Error ? error.message : 'Unknown error'}`,
    });
  };

  return (
    <div className='container mx-auto space-y-8 px-4 py-10'>
      <div className='text-center'>
        <h1 className='font-bold text-3xl'>Modlists</h1>
        <p className='mt-2 text-muted-foreground'>Manage your imported modlists</p>
      </div>

      <div className='relative mx-auto max-w-md'>
        <SearchIcon className='-translate-y-1/2 absolute top-1/2 left-3 h-4 w-4 text-muted-foreground' />
        <Input
          placeholder='Search modlists...'
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

      {filteredModlists.length === 0 ? (
        <div className='py-12 text-center'>
          <p className='text-muted-foreground'>
            {searchTerm
              ? 'No modlists found matching your search.'
              : 'No modlists found. Import a Wabbajack file to get started.'}
          </p>
        </div>
      ) : (
        <div className='grid gap-6 md:grid-cols-2 lg:grid-cols-3'>
          {filteredModlists.map(modlist => (
            <Card key={modlist.id} className='gap-2 overflow-hidden pt-0 transition-shadow hover:shadow-lg'>
              <Link to='/modlists/$id' params={{ id: modlist.id }}>
                <ModlistImage
                  modlistId={modlist.id}
                  image={modlist.image}
                  alt={modlist.name}
                  className='h-48 w-full'
                  roundedTop={true}
                />
              </Link>
              <CardHeader>
                <div className='flex items-start justify-between'>
                  <Link to='/modlists/$id' params={{ id: modlist.id }} className='flex-1'>
                    <div className='flex-1'>
                      <CardTitle className='text-lg transition-colors hover:text-blue-600'>{modlist.name}</CardTitle>
                      {modlist.author && <CardDescription className='mt-1'>by {modlist.author}</CardDescription>}
                    </div>
                  </Link>
                  <div className='flex items-center gap-2'>
                    {modlist.is_nsfw && <Badge variant='destructive'>NSFW</Badge>}
                    <AlertDialog>
                      <AlertDialogTrigger asChild>
                        <Button variant='ghost' size='sm' onClick={e => e.stopPropagation()}>
                          <Trash2 className='h-4 w-4 text-red-500' />
                        </Button>
                      </AlertDialogTrigger>
                      <AlertDialogContent>
                        <AlertDialogHeader>
                          <AlertDialogTitle>Are you sure?</AlertDialogTitle>
                          <AlertDialogDescription>
                            This action cannot be undone. This will permanently delete the modlist "{modlist.name}" and
                            all associated files.
                          </AlertDialogDescription>
                        </AlertDialogHeader>
                        <AlertDialogFooter>
                          <AlertDialogCancel>Cancel</AlertDialogCancel>
                          <AlertDialogAction onClick={() => handleDeleteModlist(modlist.id)}>Delete</AlertDialogAction>
                        </AlertDialogFooter>
                      </AlertDialogContent>
                    </AlertDialog>
                  </div>
                </div>
              </CardHeader>
              <CardContent>
                <Link to='/modlists/$id' params={{ id: modlist.id }}>
                  <div className='space-y-2'>
                    {modlist.description && (
                      <p className='line-clamp-3 text-muted-foreground text-sm'>{modlist.description}</p>
                    )}

                    <div className='flex flex-wrap gap-2 text-xs'>
                      {modlist.game_type && (
                        <Badge variant='secondary'>{pascalCaseToTitleCase(modlist.game_type)}</Badge>
                      )}
                      {modlist.version && <Badge variant='outline'>v{modlist.version}</Badge>}
                      {modlist.is_nsfw && <Badge variant='destructive'>NSFW</Badge>}
                    </div>

                    <div className='pt-2 text-muted-foreground text-xs'>
                      Imported: {new Date(modlist.created_at).toLocaleDateString()}
                    </div>
                  </div>
                </Link>
              </CardContent>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}
