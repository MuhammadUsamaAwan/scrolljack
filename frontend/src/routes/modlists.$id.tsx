import { useSuspenseQuery } from '@tanstack/react-query';
import { createFileRoute, useParams } from '@tanstack/react-router';
import { ModlistInfo } from '~/components/modlist-info';
import { queryClient } from '~/lib/query-client';
import { modListQueryOptions } from '~/lib/query-options';

export const Route = createFileRoute('/modlists/$id')({
  component: RouteComponent,
  loader: async ({ params }) => {
    const data = await queryClient.ensureQueryData(modListQueryOptions(params.id));
    return data;
  },
});

function RouteComponent() {
  const { id } = useParams({ from: '/modlists/$id' });
  const { data } = useSuspenseQuery(modListQueryOptions(id));

  return (
    <div className='container mx-auto space-y-8 px-4 py-10'>
      <ModlistInfo modlist={data} />
    </div>
  );
}
