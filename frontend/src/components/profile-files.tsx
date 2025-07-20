import { useQuery } from '@tanstack/react-query';
import { toast } from 'sonner';
import { Badge } from '~/components/ui/badge';
import { Skeleton } from '~/components/ui/skeleton';
import { profileFilesQueryOptions } from '~/lib/query-options';
import { DownloadFile } from '~/wailsjs/go/main/App';

export function ProfileFiles({ profileId }: { profileId: string }) {
  const { data, isPending } = useQuery(profileFilesQueryOptions(profileId));

  if (isPending) {
    return <Skeleton className='w-full h-8' />;
  }

  return (
    <div className='flex flex-wrap gap-4'>
      {data?.map(f => (
        <button
          key={f.id}
          type='button'
          className='cursor-pointer'
          onClick={() => {
            toast.promise(DownloadFile(f.file_path, f.name), {
              loading: 'Downloading file...',
              success: 'File downloaded successfully!',
              error: error => `Error downloading file: ${error instanceof Error ? error.message : 'Unknown error'}`,
            });
          }}
        >
          <Badge variant='outline'>{f.name}</Badge>
        </button>
      ))}
    </div>
  );
}
