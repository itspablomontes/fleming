import { createRootRoute, Outlet } from '@tanstack/react-router'

export const Route = createRootRoute({
  component: () => (
    <>
      <div className="p-2 gap-2 flex">
        {/* Placeholder for future navigation */}
      </div>
      <hr />
      <Outlet />
    </>
  ),
})
