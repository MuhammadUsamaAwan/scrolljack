import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/modlists/')({
  component: RouteComponent,
})

function RouteComponent() {
  return <div>Hello "/modlists/"!</div>
}
