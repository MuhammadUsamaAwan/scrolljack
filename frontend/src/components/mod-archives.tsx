import { useQuery } from '@tanstack/react-query';
import { Badge } from '~/components/ui/badge';
import { Spinner } from '~/components/ui/spinner';
import { modArchivesQueryOptions } from '~/lib/query-options';

export function ModArchives({ modId }: { modId: string }) {
  const { data: archives, isPending } = useQuery(modArchivesQueryOptions(modId));

  if (isPending) {
    return <Spinner />;
  }

  const gameSourceFiles = archives?.filter(a => a.type === 'GameFileSourceDownloader, Wabbajack.Lib') ?? [];
  const otherArchives = archives?.filter(a => a.type !== 'GameFileSourceDownloader, Wabbajack.Lib') ?? [];

  return (
    <>
      <div className='mb-2 text-muted-foreground text-sm'>Mod files are from {archives?.length} archive(s)</div>
      <div className='space-y-2'>
        {gameSourceFiles.length > 0 && (
          <div className='space-y-1'>
            <div className='text-muted-foreground text-sm'>
              {gameSourceFiles.length} Game Source File{gameSourceFiles.length > 1 ? 's' : ''}
            </div>
            {gameSourceFiles.map(
              a =>
                a.description && (
                  <div key={a.id} className='text-muted-foreground text-sm'>
                    {a.description}
                  </div>
                )
            )}
          </div>
        )}
        {otherArchives.map(a => (
          <div className='space-y-1' key={a.id}>
            {a.type === 'NexusDownloader, Wabbajack.Lib' ? (
              <div className='flex flex-wrap gap-2'>
                <a
                  href={`https://www.nexusmods.com/${a.nexus_game_name?.toLowerCase()}/mods/${a.nexus_mod_id}`}
                  target='_blank'
                  rel='noopener noreferrer'
                >
                  <Badge variant='outline'>Nexus Page</Badge>
                </a>
                <a
                  href={`https://www.nexusmods.com/${a.nexus_game_name?.toLowerCase()}/mods/${a.nexus_mod_id}?tab=files`}
                  target='_blank'
                  rel='noopener noreferrer'
                >
                  <Badge variant='outline'>Nexus Files</Badge>
                </a>
                <a
                  href={`https://www.nexusmods.com/${a.nexus_game_name?.toLowerCase()}/mods/${a.nexus_mod_id}?tab=files&file_id=${a.nexus_file_id}&mm==1`}
                  target='_blank'
                  rel='noopener noreferrer'
                >
                  <Badge variant='outline'>Mod Manager Download</Badge>
                </a>
                <a
                  href={`https://www.nexusmods.com/${a.nexus_game_name?.toLowerCase()}/mods/${a.nexus_mod_id}?tab=files&file_id=${a.nexus_file_id}`}
                  target='_blank'
                  rel='noopener noreferrer'
                >
                  <Badge variant='outline'>Manual Download</Badge>
                </a>
              </div>
            ) : (
              <a href={a.direct_url!} target='_blank' rel='noopener noreferrer'>
                Direct Download
              </a>
            )}
            {a.description && <div className='text-muted-foreground text-sm'>{a.description}</div>}
          </div>
        ))}
      </div>
    </>
  );
}
