import { useMutation } from "@tanstack/react-query";
import type { JSX } from "react";
import { useState } from "react";

import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Separator } from "@/components/ui/separator";

import { verifyRootOnChain } from "../api";

export function VerifyRootWidget(): JSX.Element {
	const [root, setRoot] = useState("");

	const verify = useMutation({
		mutationFn: async () => verifyRootOnChain(root.trim()),
	});

	const result = verify.data;

	return (
		<Card className="bg-white dark:bg-gray-900">
			<CardHeader className="flex flex-col gap-1">
				<CardTitle className="text-sm font-semibold text-gray-900 dark:text-gray-100">
					Verify root on-chain
				</CardTitle>
				<p className="text-xs text-muted-foreground">
					Check if a Merkle root is anchored on the configured chain.
				</p>
			</CardHeader>
			<CardContent className="space-y-3">
				<div className="flex gap-2">
					<Input
						value={root}
						onChange={(e) => setRoot(e.target.value)}
						placeholder="0x… or 64-hex root"
					/>
					<Button
						variant="outline"
						disabled={!root.trim() || verify.isPending}
						onClick={() => verify.mutate()}
					>
						{verify.isPending ? "Verifying..." : "Verify"}
					</Button>
				</div>

				{verify.isError ? (
					<p className="text-xs text-destructive">Failed to verify root.</p>
				) : null}

				{result ? (
					<>
						<Separator />
						<div className="space-y-1 text-xs text-muted-foreground">
							<p>Anchored: {result.anchored ? "true" : "false"}</p>
							<p>Timestamp: {result.timestamp}</p>
							<p>Block number: {result.blockNumber ?? "null"}</p>
							<p className="truncate">Tx hash: {result.txHash ?? "null"}</p>
						</div>
					</>
				) : null}
			</CardContent>
		</Card>
	);
}
