import { createFileRoute } from "@tanstack/react-router";
import { RootPage } from "@/features/timeline/pages/root-page";

export const Route = createFileRoute("/")({
	component: RootPage,
});
