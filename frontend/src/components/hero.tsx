import { ScrollTextIcon } from "lucide-react";

export function Hero() {
  return (
    <div className='text-center'>
      <h1 className='font-bold text-5xl flex items-center gap-1 justify-center'><ScrollTextIcon className="size-13" />Scrolljack</h1>
      <p className='mt-4 text-lg text-muted-foreground'>
        Turn your <code>.wabbajack</code> files into readable modding guides
      </p>
    </div>
  );
}
