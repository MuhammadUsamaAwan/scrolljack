import { useQuery } from '@tanstack/react-query';
import { GetModlistImageBase64 } from '~/wailsjs/go/main/App';

interface ModlistImageProps {
  modlistId: string;
  image: string;
  alt: string;
  className?: string;
}

export function ModlistImage({ modlistId, image, alt, className }: ModlistImageProps) {
  const { data, isPending } = useQuery({
    queryKey: ['modlistImage', modlistId, image],
    queryFn: async () => {
      return await GetModlistImageBase64(modlistId, image);
    },
    enabled: !!modlistId && !!image,
  });

  if (isPending) {
    return (
      <div className={`animate-pulse bg-muted rounded-t-lg ${className}`}>
        <div className='flex h-full w-full items-center justify-center text-muted-foreground text-sm'>Loading...</div>
      </div>
    );
  }

  if (!data) {
    return (
      <div className={`flex items-center justify-center bg-muted rounded-t-lg ${className}`}>
        <div className='text-muted-foreground text-sm'>No image</div>
      </div>
    );
  }

  return <img src={data} alt={alt} className={`object-cover rounded-t-lg ${className}`} />;
}
