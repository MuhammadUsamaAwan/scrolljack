import { Collapsible, CollapsibleContent, CollapsibleTrigger } from '~/components/ui/collapsible';
import { cn } from '~/lib/utils';
import { DetectFomodOptions } from '~/wailsjs/go/main/App';
import { dtos } from '~/wailsjs/go/models';
import { ModArchives } from './mod-archives';
import { ModFiles } from './mod-files';

export function Mod({ mod }: { mod: dtos.ModDTO }) {
  return (
    <Collapsible className='rounded-lg border bg-card'>
      <CollapsibleTrigger className='flex w-full cursor-pointer items-center justify-between px-4 py-2.5 after:text-muted-foreground after:text-xs after:duration-100 after:content-["⮞"] aria-expanded:after:rotate-90'>
        <span className={cn(!mod.is_active && 'text-red-500 line-through')}>
          {mod.mod_order}. {mod.name}
        </span>
      </CollapsibleTrigger>
      <CollapsibleContent className='space-y-3 px-4 pb-2.5'>
        <ModArchives modId={mod.id} />
        <Collapsible>
          <CollapsibleTrigger className='cursor-pointer text-muted-foreground text-sm underline'>
            Show/Hide Files
          </CollapsibleTrigger>
          <CollapsibleContent>
            <ModFiles modId={mod.id} />
          </CollapsibleContent>
        </Collapsible>
        <button type='button' className='underline text-sm text-muted-foreground cursor-pointer' onClick={() => DetectFomodOptions(mod.id)}>Detect Fomod Options</button>
      </CollapsibleContent>
    </Collapsible>
  );
}
