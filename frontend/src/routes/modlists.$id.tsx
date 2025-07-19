import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/modlists/$id')({
  component: RouteComponent,
})

function RouteComponent() {
  return <div>Hello "/modlists/$id"!</div>
}
