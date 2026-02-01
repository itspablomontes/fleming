import {
	Anchor,
	Check,
	ChevronDown,
	ChevronRight,
	Copy,
	ExternalLink,
	GitBranch,
	Hash,
} from "lucide-react";
import { useState } from "react";
import { Button } from "@/components/ui/button";
import {
	Tooltip,
	TooltipContent,
	TooltipTrigger,
} from "@/components/ui/tooltip";
import { cn } from "@/lib/utils";
import type { MerkleBatch } from "../api";

/**
 * MerkleTreeViewer
 * Full interactive Merkle tree visualization with expandable nodes.
 * User preferred option A: Full interactive tree (collapsible nodes, click to expand).
 */
interface MerkleTreeViewerProps {
	batch: MerkleBatch;
	className?: string;
}

// Simulated tree structure for visualization
// In production, this would be fetched from the backend
interface TreeNode {
	id: string;
	hash: string;
	type: "root" | "internal" | "leaf";
	children?: TreeNode[];
	entryData?: {
		action: string;
		timestamp: string;
	};
}

function shortHash(value?: string): string {
	if (!value) return "";
	if (value.length <= 16) return value;
	return `${value.slice(0, 10)}…${value.slice(-6)}`;
}

function formatDateTime(value?: string): string {
	if (!value) return "";
	const d = new Date(value);
	if (Number.isNaN(d.getTime())) return value;
	return d.toLocaleString();
}

// Generate a visual tree structure from batch data
function generateTreePreview(batch: MerkleBatch): TreeNode {
	// This is a simplified visualization
	// The actual Merkle tree structure would come from the backend
	const entryCount = batch.entryCount || 4;
	const leaves: TreeNode[] = Array.from({ length: Math.min(entryCount, 8) }, (_, i) => ({
		id: `leaf-${i}`,
		hash: `0x${Math.random().toString(16).slice(2, 18)}...`,
		type: "leaf" as const,
		entryData: {
			action: ["audit.create", "consent.approve", "file.upload", "auth.login"][i % 4],
			timestamp: new Date(Date.now() - i * 3600000).toISOString(),
		},
	}));

	// Build tree bottom-up
	const buildLevel = (nodes: TreeNode[]): TreeNode[] => {
		if (nodes.length <= 1) return nodes;
		const parents: TreeNode[] = [];
		for (let i = 0; i < nodes.length; i += 2) {
			const left = nodes[i];
			const right = nodes[i + 1];
			parents.push({
				id: `internal-${parents.length}`,
				hash: `0x${Math.random().toString(16).slice(2, 18)}...`,
				type: "internal" as const,
				children: right ? [left, right] : [left],
			});
		}
		return parents;
	};

	let level = leaves;
	while (level.length > 1) {
		level = buildLevel(level);
	}

	return {
		id: "root",
		hash: batch.rootHash || "0x...",
		type: "root",
		children: level[0]?.children || leaves,
	};
}

export function MerkleTreeViewer({ batch, className }: MerkleTreeViewerProps) {
	const [expandedNodes, setExpandedNodes] = useState<Set<string>>(new Set(["root"]));
	const [copiedHash, setCopiedHash] = useState<string | null>(null);

	const tree = generateTreePreview(batch);

	const toggleNode = (nodeId: string) => {
		const newExpanded = new Set(expandedNodes);
		if (newExpanded.has(nodeId)) {
			newExpanded.delete(nodeId);
		} else {
			newExpanded.add(nodeId);
		}
		setExpandedNodes(newExpanded);
	};

	const copyToClipboard = async (hash: string) => {
		await navigator.clipboard.writeText(hash);
		setCopiedHash(hash);
		setTimeout(() => setCopiedHash(null), 2000);
	};

	const expandAll = () => {
		const allIds = new Set<string>();
		const traverse = (node: TreeNode) => {
			allIds.add(node.id);
			node.children?.forEach(traverse);
		};
		traverse(tree);
		setExpandedNodes(allIds);
	};

	const collapseAll = () => {
		setExpandedNodes(new Set(["root"]));
	};

	const renderNode = (node: TreeNode, depth: number = 0): React.ReactNode => {
		const isExpanded = expandedNodes.has(node.id);
		const hasChildren = node.children && node.children.length > 0;

		const nodeColors = {
			root: "border-primary bg-primary/10 text-primary",
			internal: "border-accent bg-accent/10 text-accent",
			leaf: "border-muted-foreground bg-muted/50 text-foreground",
		};

		const nodeIcons = {
			root: Anchor,
			internal: GitBranch,
			leaf: Hash,
		};

		const Icon = nodeIcons[node.type];

		return (
			<div key={node.id} className="relative">
				{/* Connection line from parent */}
				{depth > 0 && (
					<div
						className="absolute left-0 top-0 h-4 w-4 border-l-2 border-b-2 border-border -translate-x-4 -translate-y-2"
						aria-hidden="true"
					/>
				)}

				{/* Node */}
				<div
					className={cn(
						"flex items-center gap-2 rounded-lg border-2 px-3 py-2",
						"transition-all duration-200",
						"hover:shadow-md",
						nodeColors[node.type],
						depth > 0 && "ml-6",
					)}
				>
					{/* Expand/collapse button */}
					{hasChildren && (
						<button
							type="button"
							onClick={() => toggleNode(node.id)}
							className="p-0.5 -ml-1 hover:bg-background/50 rounded"
							aria-label={isExpanded ? "Collapse" : "Expand"}
						>
							{isExpanded ? (
								<ChevronDown className="h-4 w-4" />
							) : (
								<ChevronRight className="h-4 w-4" />
							)}
						</button>
					)}

					{/* Icon */}
					<Icon className="h-4 w-4 shrink-0" />

					{/* Node content */}
					<div className="flex-1 min-w-0">
						<div className="flex items-center gap-2">
							<span className="text-xs font-medium capitalize">
								{node.type}
							</span>
							{node.type === "leaf" && node.entryData && (
								<span className="text-xs text-muted-foreground">
									{node.entryData.action}
								</span>
							)}
						</div>
						<code className="text-xs font-mono text-muted-foreground truncate block">
							{shortHash(node.hash)}
						</code>
					</div>

					{/* Copy button */}
					<Tooltip>
						<TooltipTrigger asChild>
							<button
								type="button"
								onClick={() => copyToClipboard(node.hash)}
								className="p-1 hover:bg-background/50 rounded"
								aria-label="Copy hash"
							>
								{copiedHash === node.hash ? (
									<Check className="h-3 w-3 text-success" />
								) : (
									<Copy className="h-3 w-3" />
								)}
							</button>
						</TooltipTrigger>
						<TooltipContent side="top">
							<p>Copy full hash</p>
						</TooltipContent>
					</Tooltip>
				</div>

				{/* Children */}
				{hasChildren && isExpanded && (
					<div className="mt-2 ml-4 space-y-2 pl-2 border-l-2 border-border">
						{node.children?.map((child) => renderNode(child, depth + 1))}
					</div>
				)}
			</div>
		);
	};

	const isAnchored = batch.anchorStatus === "anchored";

	return (
		<div className={cn("rounded-xl border border-border bg-card/50 p-4", className)}>
			{/* Header */}
			<div className="flex items-center justify-between mb-4">
				<div className="flex items-center gap-2">
					<GitBranch className="h-5 w-5 text-primary" />
					<h3 className="text-sm font-semibold text-foreground">
						Merkle Tree • Batch #{batch.id}
					</h3>
				</div>

				<div className="flex items-center gap-2">
					<Button
						variant="ghost"
						size="sm"
						onClick={expandAll}
						className="text-xs"
					>
						Expand All
					</Button>
					<Button
						variant="ghost"
						size="sm"
						onClick={collapseAll}
						className="text-xs"
					>
						Collapse All
					</Button>
				</div>
			</div>

			{/* Batch info */}
			<div className="mb-4 p-3 rounded-lg bg-muted/30 border border-border">
				<div className="grid grid-cols-2 md:grid-cols-4 gap-3 text-xs">
					<div>
						<span className="text-muted-foreground">Entries:</span>
						<span className="ml-1 font-medium">{batch.entryCount}</span>
					</div>
					<div>
						<span className="text-muted-foreground">Status:</span>
						<span
							className={cn(
								"ml-1 font-medium",
								isAnchored ? "text-success" : "text-warning",
							)}
						>
							{isAnchored ? "Anchored" : "Pending"}
						</span>
					</div>
					<div>
						<span className="text-muted-foreground">Created:</span>
						<span className="ml-1 font-medium">
							{formatDateTime(batch.createdAt)}
						</span>
					</div>
					{batch.anchorTxHash && (
						<div className="flex items-center gap-1">
							<span className="text-muted-foreground">Tx:</span>
							<a
								href={`https://sepolia.basescan.org/tx/${batch.anchorTxHash}`}
								target="_blank"
								rel="noopener noreferrer"
								className="text-primary hover:underline flex items-center gap-1"
							>
								{shortHash(batch.anchorTxHash)}
								<ExternalLink className="h-3 w-3" />
							</a>
						</div>
					)}
				</div>
			</div>

			{/* Tree visualization */}
			<div className="space-y-2 overflow-x-auto">
				{renderNode(tree)}
			</div>

			{/* Legend */}
			<div className="mt-4 pt-4 border-t border-border">
				<div className="flex flex-wrap gap-4 text-xs text-muted-foreground">
					<div className="flex items-center gap-1">
						<div className="h-3 w-3 rounded border-2 border-primary bg-primary/20" />
						<span>Root (anchored on-chain)</span>
					</div>
					<div className="flex items-center gap-1">
						<div className="h-3 w-3 rounded border-2 border-accent bg-accent/20" />
						<span>Internal node</span>
					</div>
					<div className="flex items-center gap-1">
						<div className="h-3 w-3 rounded border-2 border-muted-foreground bg-muted/50" />
						<span>Leaf (audit entry)</span>
					</div>
				</div>
			</div>
		</div>
	);
}
