import { useQuery, useQueryClient } from "@tanstack/react-query";
import { Link } from "@tanstack/react-router";
import { ArrowLeft, Plus, Shield } from "lucide-react";
import type { JSX } from "react";
import { useMemo, useState } from "react";
import { Button } from "@/components/ui/button";
import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogHeader,
	DialogTitle,
	DialogTrigger,
} from "@/components/ui/dialog";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { queryAuditEntries } from "@/features/audit/api";
import { AuditLogTable } from "@/features/audit/components/audit-log-table";
import { AuditTargetType } from "@/features/audit/types";

import { ConsentDashboard } from "../components/consent-dashboard";
import { ConsentRequestWizard } from "../components/consent-request-wizard";
import { ConsentStatsBar } from "../components/consent-stats-bar";

import { useActiveConsentGrants, useMyConsentGrants } from "../hooks/use-consent-grants";

const auditPreviewLimit = 5;

export function ConsentPage(): JSX.Element {
	const [wizardOpen, setWizardOpen] = useState(false);
	const [activeFilter, setActiveFilter] = useState<"pending" | "active" | "expired" | null>(null);
	const queryClient = useQueryClient();

	const myGrantsQuery = useMyConsentGrants();
	const activeGrantsQuery = useActiveConsentGrants();

	const auditQuery = useQuery({
		queryKey: ["audit-logs", "consent-preview"],
		queryFn: () =>
			queryAuditEntries({
				resourceType: AuditTargetType.Consent,
				limit: auditPreviewLimit,
			}),
	});

	// Calculate counts for summary cards
	const counts = useMemo(() => {
		const myGrants = myGrantsQuery.data ?? [];
		const isExpired = (grant: { expiresAt?: Date | null }) =>
			Boolean(grant.expiresAt && grant.expiresAt.getTime() <= Date.now());

		const pending = myGrants.filter((g) => g.state === "requested").length;
		const active = myGrants.filter(
			(g) => g.state === "approved" && !isExpired(g),
		).length;
		const expired = myGrants.filter((g) => {
			if (g.state === "requested") return false;
			if (g.state === "approved") return isExpired(g);
			return true; // revoked, denied, etc.
		}).length;

		return { pending, active, expired };
	}, [myGrantsQuery.data]);

	const handleWizardSuccess = () => {
		setWizardOpen(false);
		void queryClient.invalidateQueries({ queryKey: ["consent"] });
		void auditQuery.refetch();
		void myGrantsQuery.refetch();
		void activeGrantsQuery.refetch();
	};

	const handleFilterClick = (filter: "pending" | "active" | "expired") => {
		setActiveFilter((prev) => (prev === filter ? null : filter));
	};

	const isLoading = auditQuery.isLoading || myGrantsQuery.isLoading;

	return (
		<div className="min-h-screen bg-background px-4 py-6 md:px-8 md:py-10">
			<div className="mx-auto flex w-full max-w-5xl flex-col gap-6">
				{/* Header */}
				<div className="flex flex-col gap-4">
					<Button
						variant="ghost"
						size="sm"
						asChild
						className="w-fit gap-2 text-muted-foreground hover:text-foreground"
						aria-label="Back to timeline"
					>
						<Link to="/">
							<ArrowLeft className="h-4 w-4" />
							Back to Timeline
						</Link>
					</Button>

					<div className="flex items-start justify-between gap-4">
						<div className="flex items-start gap-4">
							<div className="flex h-12 w-12 items-center justify-center rounded-xl bg-primary/10 shrink-0">
								<Shield className="h-6 w-6 text-primary" />
							</div>
							<div className="flex flex-col gap-1">
								<h1 className="text-2xl font-bold tracking-tight text-foreground">
									Access Control
								</h1>
								<p className="text-muted-foreground">
									Control who can view or update your medical timeline.
								</p>
							</div>
						</div>

						{/* Request Access Button - Opens Modal */}
						<Dialog open={wizardOpen} onOpenChange={setWizardOpen}>
							<DialogTrigger asChild>
								<Button className="gap-2 shrink-0">
									<Plus className="h-4 w-4" />
									Request Access
								</Button>
							</DialogTrigger>
							<DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
								<DialogHeader>
									<DialogTitle>Request Access</DialogTitle>
									<DialogDescription>
										Send a consent request to a grantor to access their timeline.
									</DialogDescription>
								</DialogHeader>
								<ConsentRequestWizard onSuccess={handleWizardSuccess} />
							</DialogContent>
						</Dialog>
					</div>
				</div>

				{/* Main Content Tabs */}
				<Tabs defaultValue="my-data" className="w-full">
					<TabsList className="grid w-full max-w-[400px] grid-cols-2">
						<TabsTrigger value="my-data">My Data</TabsTrigger>
						<TabsTrigger value="external-access">External Access</TabsTrigger>
					</TabsList>
					
					<TabsContent value="my-data" className="mt-6 space-y-6">
						{/* Summary Stats Bar */}
						<ConsentStatsBar
							pendingCount={counts.pending}
							activeCount={counts.active}
							expiredCount={counts.expired}
							onFilterClick={handleFilterClick}
							activeFilter={activeFilter}
							isLoading={isLoading}
						/>

						<ConsentDashboard
							mode="grantor"
							defaultTab={
								activeFilter === "expired"
									? "history"
									: activeFilter ?? undefined
							}
						/>
					</TabsContent>

					<TabsContent value="external-access" className="mt-6">
						<ConsentDashboard
							mode="grantee"
							defaultTab="active"
						/>
					</TabsContent>
				</Tabs>

				{/* Recent Consent Activity */}
				<div className="rounded-xl border border-border bg-card p-4 space-y-4">
					<div className="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
						<div className="space-y-1">
							<h2 className="text-base font-semibold text-foreground">
								Recent consent activity
							</h2>
							<p className="text-sm text-muted-foreground">
								Last {auditPreviewLimit} consent-related actions.
							</p>
						</div>
						<Button variant="outline" size="sm" asChild>
							<Link to="/audit">View full audit log</Link>
						</Button>
					</div>
					<AuditLogTable
						entries={auditQuery.data?.entries ?? []}
						isLoading={auditQuery.isLoading}
					/>
				</div>
			</div>
		</div>
	);
}
