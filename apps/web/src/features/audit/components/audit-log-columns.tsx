"use client";

import type { ColumnDef } from "@tanstack/react-table";
import { 
    ArrowUpDown, 
    Copy, 
    ExternalLink, 
    FileText, 
    Key, 
    MoreHorizontal, 
    ShieldCheck, 
    User 
} from "lucide-react";
import { toast } from "sonner";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import { cn } from "@/lib/utils";
import type { AuditLogEntry } from "../types";
import { formatAuditAction, formatAuditResourceType } from "../types";

// Helper for copy
const copyToClipboard = (text: string, label: string) => {
    navigator.clipboard.writeText(text);
    toast.success(`${label} copied to clipboard`);
};

// Helper for address
const shortenAddress = (address: string) => {
    if (!address || address.length < 10) return address;
    return `${address.substring(0, 6)}...${address.substring(address.length - 4)}`;
};

export const columns: ColumnDef<AuditLogEntry>[] = [
    {
        accessorKey: "timestamp",
        header: ({ column }) => {
            return (
                <Button
                    variant="ghost"
                    onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
                    className="-ml-4 whitespace-nowrap"
                >
                    Timestamp
                    <ArrowUpDown className="ml-2 h-4 w-4" />
                </Button>
            );
        },
        cell: ({ row }) => {
            const date = new Date(row.getValue("timestamp"));
            return (
                <div className="flex flex-col text-sm">
                    <span className="font-medium text-foreground">
                        {date.toLocaleDateString(undefined, { month: "short", day: "numeric", year: "numeric" })}
                    </span>
                    <span className="text-xs text-muted-foreground">
                        {date.toLocaleTimeString(undefined, { hour: "2-digit", minute: "2-digit", second: "2-digit" })}
                    </span>
                </div>
            )
        },
    },
    {
        accessorKey: "resourceType",
        header: "Type",
        cell: ({ row }) => {
            const type = row.getValue("resourceType") as string;
            const icon = getIconForType(type);
            const label = formatAuditResourceType(type);
            
            return (
                <div className="flex items-center gap-2">
                    <div className="flex h-8 w-8 items-center justify-center rounded-md border border-border bg-muted/50">
                        {icon}
                    </div>
                    <span className="font-medium text-sm text-muted-foreground">{label}</span>
                </div>
            );
        },
    },
    {
        accessorKey: "action",
        header: "Action",
        cell: ({ row }) => {
            const action = row.getValue("action") as string;
            const label = formatAuditAction(action);
            const styles = getActionStyles(action);
            
            return (
                <Badge variant="outline" className={cn("font-normal capitalize", styles)}>
                    {label}
                </Badge>
            );
        },
    },
    {
        accessorKey: "resourceId",
        header: "Resource ID",
        cell: ({ row }) => {
            const id = row.getValue("resourceId") as string;
            return (
                <div className="flex items-center gap-2 group">
                    <code className="rounded bg-muted px-[0.3rem] py-[0.2rem] font-mono text-xs text-muted-foreground group-hover:text-foreground transition-colors">
                        {id.substring(0, 8)}...
                    </code>
                    <Button
                        variant="ghost"
                        size="icon"
                        className="h-6 w-6 opacity-0 group-hover:opacity-100 transition-opacity"
                        onClick={() => copyToClipboard(id, "Resource ID")}
                    >
                        <Copy className="h-3 w-3" />
                         <span className="sr-only">Copy ID</span>
                    </Button>
                </div>
            );
        },
    },
    {
        accessorKey: "actor",
        header: "Actor",
        cell: ({ row }) => {
            const actor = row.getValue("actor") as string;
            return (
                <div className="flex items-center gap-2">
                     <div className="flex h-6 w-6 items-center justify-center rounded-full bg-primary/10">
                        <User className="h-3 w-3 text-primary" />
                     </div>
                     <span className="font-mono text-xs text-muted-foreground">
                        {shortenAddress(actor)}
                     </span>
                </div>
            );
        },
    },
    {
        id: "integrity",
        header: () => <div className="text-right">Integrity</div>,
        cell: ({ row }) => {
            const hash = row.original.hash;
            if (!hash) return <div className="text-right text-muted-foreground">-</div>;
            
            return (
                <div className="flex justify-end">
                    <TooltipProvider delayDuration={0}>
                        <Tooltip>
                            <TooltipTrigger asChild>
                                <Button 
                                    variant="ghost" 
                                    size="icon" 
                                    className="h-8 w-8 text-green-600 hover:text-green-700 hover:bg-green-50 dark:text-green-500 dark:hover:bg-green-950/50"
                                    onClick={() => copyToClipboard(hash, "Integrity Hash")}
                                >
                                    <ShieldCheck className="h-4 w-4" />
                                    <span className="sr-only">Verified</span>
                                </Button>
                            </TooltipTrigger>
                            <TooltipContent side="left" className="font-mono text-xs max-w-[300px] break-all">
                                <div className="space-y-1">
                                    <p className="font-semibold text-green-600 flex items-center gap-2">
                                        <ShieldCheck className="h-3 w-3" /> Cryptographically Verified
                                    </p>
                                    <p className="text-muted-foreground">{hash}</p>
                                    <p className="text-[10px] text-muted-foreground/50 pt-1">
                                        Click icon to copy full hash
                                    </p>
                                </div>
                            </TooltipContent>
                        </Tooltip>
                    </TooltipProvider>
                </div>
            );
        },
    },
    {
        id: "actions",
        enableHiding: false,
        cell: ({ row }) => {
            const entry = row.original;
            
            return (
                <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                        <Button variant="ghost" className="h-8 w-8 p-0">
                            <span className="sr-only">Open menu</span>
                            <MoreHorizontal className="h-4 w-4" />
                        </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="end">
                        <DropdownMenuLabel>Actions</DropdownMenuLabel>
                        <DropdownMenuItem onClick={() => copyToClipboard(entry.id, "Entry ID")}>
                            Copy Entry ID
                        </DropdownMenuItem>
                        <DropdownMenuItem onClick={() => copyToClipboard(JSON.stringify(entry), "JSON")}>
                            Copy JSON
                        </DropdownMenuItem>
                    </DropdownMenuContent>
                </DropdownMenu>
            )
        }
    }
];

// Icons helper
function getIconForType(type: string) {
	switch (type) {
		case "file":
			return <FileText className="h-4 w-4 text-blue-500" />;
		case "consent":
			return <ShieldCheck className="h-4 w-4 text-green-500" />;
		case "auth":
			return <Key className="h-4 w-4 text-orange-500" />;
		default:
			return <ExternalLink className="h-4 w-4 text-muted-foreground" />;
	}
}

// Styles helper
function getActionStyles(action: string): string {
	if (action.includes("create") || action.includes("upload") || action.includes("resume") || action.includes("approve"))
		return "border-green-500/20 text-green-700 bg-green-50 dark:bg-green-500/10 dark:text-green-400";
	if (action.includes("delete") || action.includes("revoke") || action.includes("deny") || action.includes("suspend"))
		return "border-red-500/20 text-red-700 bg-red-50 dark:bg-red-500/10 dark:text-red-400";
	if (action.includes("update") || action.includes("share"))
		return "border-blue-500/20 text-blue-700 bg-blue-50 dark:bg-blue-500/10 dark:text-blue-400";
    if (action.includes("request"))
        return "border-amber-500/20 text-amber-700 bg-amber-50 dark:bg-amber-500/10 dark:text-amber-400";
	return "border-border text-foreground bg-muted";
}
