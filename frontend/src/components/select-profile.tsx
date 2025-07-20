import { Label } from '~/components/ui/label';
import { Select, SelectContent, SelectGroup, SelectItem, SelectTrigger, SelectValue } from '~/components/ui/select';
import { models } from '~/wailsjs/go/models';

export function SelectProfile({
  profiles,
  selectedProfile,
  setSelectedProfile,
}: {
  profiles: models.Profile[];
  selectedProfile: string;
  setSelectedProfile: (id: string) => void;
}) {
  return (
    <div className='space-y-2'>
      <Label>Select Profile</Label>
      <Select value={selectedProfile} onValueChange={setSelectedProfile}>
        <SelectTrigger className='w-full'>
          <SelectValue placeholder='Select a Profile' />
        </SelectTrigger>
        <SelectContent>
          <SelectGroup>
            {profiles?.map(p => (
              <SelectItem key={p.id} value={p.id} onClick={() => setSelectedProfile(p.id)}>
                {p.name}
              </SelectItem>
            ))}
          </SelectGroup>
        </SelectContent>
      </Select>
    </div>
  );
}
