import { useQuery } from '@tanstack/react-query';
import { Fragment, useState } from 'react';
import { toast } from 'sonner';
import { FileDiff } from '~/components/file-diff';
import { Spinner } from '~/components/ui/spinner';
import { modFilesQueryOptions } from '~/lib/query-options';
import { base64ToUint8Array, formatSize, uint8ArrayToString } from '~/lib/utils';
import { ApplyBinaryPatch, DownloadFile } from '~/wailsjs/go/main/App';
import { dtos } from '~/wailsjs/go/models';

export function ModFiles({ modId }: { modId: string }) {
  const { data: files, isPending } = useQuery(modFilesQueryOptions(modId));
  const [diffFileId, setDiffFileId] = useState<string | null>(null);
  const [originalFileContent, setOriginalFileContent] = useState('');
  const [patchedFileContent, setPatchedFileContent] = useState('');

  if (isPending) {
    return <Spinner />;
  }

  async function handleBinaryPatchClick(f: dtos.ModFileDTO) {
    const { original, patched } = await ApplyBinaryPatch(f.patch_file_path!, f.path.split('\\').pop()!);
    const originalBytes = base64ToUint8Array(original);
    const patchedBytes = base64ToUint8Array(patched);
    setOriginalFileContent(uint8ArrayToString(originalBytes));
    setPatchedFileContent(uint8ArrayToString(patchedBytes));
    setDiffFileId(f.id);
  }

  return files?.map(f => (
    <Fragment key={f.id}>
      <div className='flex items-center gap-2 font-mono text-muted-foreground text-xs'>
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
              toast.promise(handleBinaryPatchClick(f), {
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
      {diffFileId === f.id && originalFileContent !== '' && patchedFileContent !== '' && (
        <div className='text-xs border rounded-xl p-4 my-2'>
          <h3 className='font-semibold mb-2'>
            File Diff{' '}
            <button type='button' className='cursor-pointer underline' onClick={() => setDiffFileId(null)}>
              Clear
            </button>
          </h3>
          <FileDiff original={originalFileContent} patched={patchedFileContent} />
        </div>
      )}
    </Fragment>
  ));
}
