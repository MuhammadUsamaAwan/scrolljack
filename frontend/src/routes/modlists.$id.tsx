import { useSuspenseQuery } from '@tanstack/react-query';
import { createFileRoute, useParams } from '@tanstack/react-router';
import { useState } from 'react';
import { ModlistInfo } from '~/components/modlist-info';
import { ProfileFiles } from '~/components/profile-files';
import { ProfileMods } from '~/components/profile-mods';
import { SelectProfile } from '~/components/select-profile';
import { queryClient } from '~/lib/query-client';
import { modListQueryOptions, profilesQueryOptions } from '~/lib/query-options';

export const Route = createFileRoute('/modlists/$id')({
  component: RouteComponent,
  loader: async ({ params }) => {
    const [modlist, profiles] = await Promise.all([
      queryClient.ensureQueryData(modListQueryOptions(params.id)),
      queryClient.ensureQueryData(profilesQueryOptions(params.id)),
    ]);
    return { modlist, profiles };
  },
});

function RouteComponent() {
  const { id } = useParams({ from: '/modlists/$id' });
  const { data: modlist } = useSuspenseQuery(modListQueryOptions(id));
  const { data: profiles } = useSuspenseQuery(profilesQueryOptions(id));
  const [selectedProfile, setSelectedProfile] = useState(profiles[0].id);

  return (
    <div className='container mx-auto space-y-8 px-4 py-10'>
      <ModlistInfo modlist={modlist} />
      <SelectProfile profiles={profiles} selectedProfile={selectedProfile} setSelectedProfile={setSelectedProfile} />
      <ProfileFiles profileId={selectedProfile} />
      <ProfileMods profileId={selectedProfile} />
    </div>
  );
}
