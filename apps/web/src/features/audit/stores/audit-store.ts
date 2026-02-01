import { create } from "zustand";
import { persist } from "zustand/middleware";
import type { AuditAction, AuditTargetType } from "../types";

export interface AuditFilterState {
	actor: string;
	resourceId: string;
	resourceType: AuditTargetType | "";
	action: AuditAction | "";
	startTime: string; // ISO date string
	endTime: string; // ISO date string
}

interface AuditStore {
	// Filters
	filters: AuditFilterState;
	setFilters: (filters: Partial<AuditFilterState>) => void;
	resetFilters: () => void;

	// Pagination
	page: number;
	pageSize: number;
	setPage: (page: number) => void;
	setPageSize: (pageSize: number) => void;
}

const defaultFilters: AuditFilterState = {
	actor: "",
	resourceId: "",
	resourceType: "",
	action: "",
	startTime: "",
	endTime: "",
};

export const useAuditStore = create<AuditStore>()(
	persist(
		(set) => ({
			filters: defaultFilters,
			page: 1,
			pageSize: 10,

			setFilters: (newFilters) =>
				set((state) => ({
					filters: { ...state.filters, ...newFilters },
					page: 1, // Reset page on filter change
				})),

			resetFilters: () =>
				set(() => ({
					filters: defaultFilters,
					page: 1,
				})),

			setPage: (page) => set(() => ({ page })),
			setPageSize: (pageSize) => set(() => ({ pageSize, page: 1 })),
		}),
		{
			name: "audit-filters-storage",
			partialize: (state) => ({ filters: state.filters, pageSize: state.pageSize }), // Don't persist current page
		},
	),
);
