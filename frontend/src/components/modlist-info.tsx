import { useMutation } from '@tanstack/react-query';
import { Link, useNavigate } from '@tanstack/react-router';
import { ArrowLeft, Trash2 } from 'lucide-react';
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
import { queryClient } from '~/lib/query-client';
import { modListQueryOptions, modListsQueryOptions } from '~/lib/query-options';
import { pascalCaseToTitleCase } from '~/lib/utils';
import { DeleteModlist } from '~/wailsjs/go/main/App';
import { dtos } from '~/wailsjs/go/models';

export function ModlistInfo({ modlist }: { modlist: dtos.ModlistDTO }) {
  const navigate = useNavigate();
  const { mutateAsync } = useMutation({
    mutationFn: DeleteModlist,
    meta: {
      invalidateQueries: [modListsQueryOptions.queryKey, modListQueryOptions(modlist.id).queryKey],
    },
    onSuccess: () => {
      navigate({ to: '/modlists' });
    },
  });

  const deleteModlist = async () => {
    await mutateAsync(modlist.id);
    queryClient.invalidateQueries({ queryKey: modListsQueryOptions.queryKey });
    queryClient.invalidateQueries({ queryKey: modListQueryOptions(modlist.id).queryKey });
  };

  const handleDeleteModlist = async () => {
    toast.promise(deleteModlist(), {
      loading: 'Deleting modlist this may take a while...',
      success: () => 'Modlist deleted successfully',
      error: error => `Failed to delete modlist: ${error instanceof Error ? error.message : 'Unknown error'}`,
    });
  };

  return (
    <>
      <div className='flex items-center justify-between'>
        <Button variant='ghost' size='sm' asChild>
          <Link to='/modlists'>
            <ArrowLeft className='h-4 w-4' />
            Back to Modlists
          </Link>
        </Button>

        <AlertDialog>
          <AlertDialogTrigger asChild>
            <Button variant='destructive' size='sm'>
              <Trash2 className='h-4 w-4' />
              Delete Modlist
            </Button>
          </AlertDialogTrigger>
          <AlertDialogContent>
            <AlertDialogHeader>
              <AlertDialogTitle>Are you sure?</AlertDialogTitle>
              <AlertDialogDescription>
                This action cannot be undone. This will permanently delete the modlist "{modlist.name}" and all
                associated files.
              </AlertDialogDescription>
            </AlertDialogHeader>
            <AlertDialogFooter>
              <AlertDialogCancel>Cancel</AlertDialogCancel>
              <AlertDialogAction onClick={handleDeleteModlist}>Delete</AlertDialogAction>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialog>
      </div>

      <Card className='gap-4 overflow-hidden pt-0'>
        <ModlistImage modlistId={modlist.id} image={modlist.image} alt={modlist.name} className='h-64 w-full' />
        <CardHeader>
          <div className='flex items-start justify-between'>
            <div className='flex-1'>
              <CardTitle className='text-2xl'>{modlist.name}</CardTitle>
              {modlist.author && <CardDescription className='mt-1 text-lg'>by {modlist.author}</CardDescription>}
            </div>
            {modlist.is_nsfw && <Badge variant='destructive'>NSFW</Badge>}
          </div>
        </CardHeader>
        <CardContent>
          <div className='space-y-4'>
            {modlist.description && <p className='text-muted-foreground'>{modlist.description}</p>}

            <div className='grid gap-4 md:grid-cols-2'>
              <div>
                <h3 className='mb-2 font-semibold'>Details</h3>
                <div className='space-y-2'>
                  {modlist.game_type && (
                    <div className='flex items-center gap-2'>
                      <span className='font-medium text-sm'>Game:</span>
                      <Badge variant='secondary'>{pascalCaseToTitleCase(modlist.game_type)}</Badge>
                    </div>
                  )}
                  {modlist.version && (
                    <div className='flex items-center gap-2'>
                      <span className='font-medium text-sm'>Version:</span>
                      <Badge variant='outline'>v{modlist.version}</Badge>
                    </div>
                  )}
                  <div className='flex items-center gap-2'>
                    <span className='font-medium text-sm'>Imported:</span>
                    <span className='text-muted-foreground text-sm'>
                      {new Date(modlist.created_at).toLocaleDateString()}
                    </span>
                  </div>
                </div>
              </div>

              <div>
                <h3 className='mb-2 font-semibold'>Links</h3>
                <div className='space-y-2'>
                  {modlist.website && (
                    <div>
                      <span className='font-medium text-sm'>Website:</span>
                      <a
                        href={modlist.website}
                        className='ml-2 text-blue-600 text-sm hover:underline'
                        target='_blank'
                        rel='noopener noreferrer'
                      >
                        {modlist.website}
                      </a>
                    </div>
                  )}
                  {modlist.readme && (
                    <div>
                      <span className='font-medium text-sm'>Readme:</span>
                      <a
                        href={modlist.readme}
                        className='ml-2 text-blue-600 text-sm hover:underline'
                        target='_blank'
                        rel='noopener noreferrer'
                      >
                        {modlist.readme}
                      </a>
                    </div>
                  )}
                </div>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>
    </>
  );
}
