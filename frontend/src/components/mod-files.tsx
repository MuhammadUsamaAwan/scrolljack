import { useQuery } from '@tanstack/react-query';
import { toast } from 'sonner';
import { Spinner } from '~/components/ui/spinner';
import { modFilesQueryOptions } from '~/lib/query-options';
import { formatSize } from '~/lib/utils';
import { ApplyBinaryPatch, DownloadFile } from '~/wailsjs/go/main/App';

export function ModFiles({ modId }: { modId: string }) {
  const { data: files, isPending } = useQuery(modFilesQueryOptions(modId));

  if (isPending) {
    return <Spinner />;
  }

  return files?.map(f => (
    <div key={f.id} className='flex items-center gap-2 font-mono text-muted-foreground text-xs'>
      {f.path}
      {(f.type === 'InlineFile' || f.type === 'RemappedInlineFile') && (
        <button
          onClick={() => {
            toast.promise(DownloadFile(f.source_file_path!, f.path.split('\\').pop()!), {
              loading: 'Downloading file...',
              success: 'File downloaded successfully!',
              error: error => `Error downloading file: ${error instanceof Error ? error.message : 'Unknown error'}`,
            });
          }}
          type='button'
          className='cursor-pointer underline'
        >
          Download
        </button>
      )}
      {f.patch_file_path && (
        <button
          onClick={() => {
            toast.promise(ApplyBinaryPatch(f.patch_file_path!, f.path.split('\\').pop()!), {
              loading: 'Applying patch...',
              success: 'Patch applied successfully and saved to the downloads folder!',
              error: error => `Error applying patch: ${error instanceof Error ? error.message : 'Unknown error'}`,
            });
          }}
          type='button'
          className='cursor-pointer underline'
        >
          Apply Patch
        </button>
      )}{' '}
      ({formatSize(f.size)}) ({f.type})
    </div>
  ));
}
