import { queryOptions } from '@tanstack/react-query';
import { GetModlistById, GetModlists } from '~/wailsjs/go/main/App';

export const modListsQueryOptions = queryOptions({
  queryKey: ['modlists'],
  queryFn: async () => {
    return await GetModlists();
  },
});

export const modListQueryOptions = (id: string) =>
  queryOptions({
    queryKey: ['modlists', id],
    queryFn: async () => {
      return await GetModlistById(id);
    },
  });
