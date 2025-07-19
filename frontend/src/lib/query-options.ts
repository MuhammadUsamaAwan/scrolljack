import { queryOptions } from '@tanstack/react-query';
import { GetModlists } from '~/wailsjs/go/main/App';

export const modListQueryOptions = queryOptions({
  queryKey: ['modlists'],
  queryFn: async () => {
    return await GetModlists();
  },
});
