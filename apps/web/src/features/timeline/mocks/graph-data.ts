/**
 * Mock Graph Data
 *
 * Realistic patient journey for development and testing.
 * Scenario: Annual checkup leads to diabetes diagnosis and treatment plan.
 *
 * =============================================================================
 * TODO: DELETE THIS ENTIRE FILE when connecting to real API
 *
 * This file provides mock data for development only.
 * When the backend API is connected:
 * 1. Delete this entire mocks/ directory
 * 2. Update get-graph.ts to use real API (set USE_MOCK_DATA = false)
 * 3. Update timeline-graph-page.tsx to use useGraphData hook
 * =============================================================================
 */

import type { EventEdge, GraphData, TimelineEvent } from "../types";

// =============================================================================
// MOCK EVENTS
// =============================================================================

export const MOCK_EVENTS: TimelineEvent[] = [
	{
		id: "evt-early-checkup",
		patientId: "0x742d35Cc6634C0532925a3b844Bc9e7595f",
		type: "consultation",
		title: "Initial Baseline Checkup",
		description:
			"Routine checkup before any chronic symptoms. All labs normal.",
		provider: "Dr. Sarah Chen",
		timestamp: "2023-03-10T10:00:00Z",
		isEncrypted: false,
		metadata: { facility: "Main Street Medical Center" },
		createdAt: "2023-03-10T10:30:00Z",
		updatedAt: "2023-03-10T10:30:00Z",
	},
	{
		id: "evt-consultation-1",
		patientId: "0x742d35Cc6634C0532925a3b844Bc9e7595f",
		type: "consultation",
		title: "Annual Checkup - Initial Symptoms",
		description:
			"Routine physical examination. Patient reports fatigue and increased thirst. Ordered CMP.",
		provider: "Dr. Sarah Chen",
		timestamp: "2024-01-15T09:00:00Z",
		isEncrypted: false,
		metadata: {
			facility: "Main Street Medical Center",
			visitType: "preventive",
		},
		createdAt: "2024-01-15T09:30:00Z",
		updatedAt: "2024-01-15T09:30:00Z",
	},
	{
		id: "evt-lab-1",
		patientId: "0x742d35Cc6634C0532925a3b844Bc9e7595f",
		type: "lab_result",
		title: "Comprehensive Metabolic Panel",
		description: "Fasting blood glucose: 186 mg/dL (High). HbA1c: 7.8% (High).",
		provider: "LabCorp",
		timestamp: "2024-01-16T14:00:00Z",
		isEncrypted: false,
		metadata: { glucoseLevel: 186, hba1c: 7.8 },
		createdAt: "2024-01-16T15:00:00Z",
		updatedAt: "2024-01-16T15:00:00Z",
	},
	{
		id: "evt-diagnosis-1",
		patientId: "0x742d35Cc6634C0532925a3b844Bc9e7595f",
		type: "diagnosis",
		title: "Type 2 Diabetes Mellitus",
		description:
			"Diagnosis confirmed based on glucose and HbA1c levels. ICD-10: E11.9",
		provider: "Dr. Sarah Chen",
		timestamp: "2024-01-20T10:00:00Z",
		isEncrypted: false,
		metadata: { icd10: "E11.9", severity: "moderate" },
		createdAt: "2024-01-20T10:30:00Z",
		updatedAt: "2024-01-20T10:30:00Z",
	},
	{
		id: "evt-prescription-1",
		patientId: "0x742d35Cc6634C0532925a3b844Bc9e7595f",
		type: "prescription",
		title: "Metformin 500mg",
		description: "Start daily dosage with dinner.",
		provider: "Dr. Sarah Chen",
		timestamp: "2024-01-20T10:15:00Z",
		isEncrypted: false,
		metadata: { medication: "Metformin", dosage: "500mg" },
		createdAt: "2024-01-20T10:30:00Z",
		updatedAt: "2024-01-20T10:30:00Z",
	},
	{
		id: "evt-consultation-2",
		patientId: "0x742d35Cc6634C0532925a3b844Bc9e7595f",
		type: "consultation",
		title: "Endocrinology Visit",
		description:
			"Discussed blood glucose monitoring and dietary modifications.",
		provider: "Dr. Michael Park",
		timestamp: "2024-02-05T14:00:00Z",
		isEncrypted: false,
		metadata: { specialty: "Endocrinology" },
		createdAt: "2024-02-05T15:00:00Z",
		updatedAt: "2024-02-05T15:00:00Z",
	},
	{
		id: "evt-lab-2",
		patientId: "0x742d35Cc6634C0532925a3b844Bc9e7595f",
		type: "lab_result",
		title: "6-Month HbA1c Review",
		description: "HbA1c: 6.8% (Significant improvement).",
		provider: "LabCorp",
		timestamp: "2024-08-10T09:00:00Z",
		isEncrypted: false,
		metadata: { hba1c: 6.8 },
		createdAt: "2024-08-10T10:00:00Z",
		updatedAt: "2024-08-10T10:00:00Z",
	},
	{
		id: "evt-procedure-1",
		patientId: "0x742d35Cc6634C0532925a3b844Bc9e7595f",
		type: "procedure",
		title: "Ophthalmology Screening",
		description: "Diabetic retinopathy screening. No issues detected.",
		provider: "Dr. Elena Rossi",
		timestamp: "2025-01-22T11:00:00Z",
		isEncrypted: false,
		metadata: { results: "clear" },
		createdAt: "2025-01-22T12:00:00Z",
		updatedAt: "2025-01-22T12:00:00Z",
	},
	{
		id: "evt-annual-2026",
		patientId: "0x742d35Cc6634C0532925a3b844Bc9e7595f",
		type: "consultation",
		title: "2026 Health Maintenance",
		description:
			"Long-term diabetes management review. HbA1c remains stable at 6.6%.",
		provider: "Dr. Sarah Chen",
		timestamp: "2026-01-15T10:00:00Z",
		isEncrypted: false,
		metadata: { status: "controlled" },
		createdAt: "2026-01-15T10:30:00Z",
		updatedAt: "2026-01-15T10:30:00Z",
	},
	{
		id: "evt-future-scan",
		patientId: "0x742d35Cc6634C0532925a3b844Bc9e7595f",
		type: "imaging",
		title: "Planned Cardiac Assessment",
		description: "Scheduled preventative scan.",
		provider: "City Imaging",
		timestamp: "2026-06-20T14:00:00Z",
		isEncrypted: false,
		metadata: { planned: true },
		createdAt: "2025-12-01T09:00:00Z",
		updatedAt: "2025-12-01T09:00:00Z",
	},
];

// =============================================================================
// MOCK EDGES
// =============================================================================

export const MOCK_EDGES: EventEdge[] = [
	{
		id: "edge-1",
		fromEventId: "evt-early-checkup",
		toEventId: "evt-consultation-1",
		relationshipType: "lead_to",
		createdAt: "2024-01-15T09:30:00Z",
	},
	{
		id: "edge-2",
		fromEventId: "evt-consultation-1",
		toEventId: "evt-lab-1",
		relationshipType: "resulted_in",
		createdAt: "2024-01-15T09:30:00Z",
	},
	{
		id: "edge-3",
		fromEventId: "evt-lab-1",
		toEventId: "evt-diagnosis-1",
		relationshipType: "supports",
		createdAt: "2024-01-20T10:30:00Z",
	},
	{
		id: "edge-4",
		fromEventId: "evt-diagnosis-1",
		toEventId: "evt-prescription-1",
		relationshipType: "lead_to",
		createdAt: "2024-01-20T10:30:00Z",
	},
	{
		id: "edge-5",
		fromEventId: "evt-prescription-1",
		toEventId: "evt-consultation-2",
		relationshipType: "resulted_in",
		createdAt: "2024-02-05T15:00:00Z",
	},
	{
		id: "edge-6",
		fromEventId: "evt-consultation-2",
		toEventId: "evt-lab-2",
		relationshipType: "resulted_in",
		createdAt: "2024-08-10T10:00:00Z",
	},
	{
		id: "edge-7",
		fromEventId: "evt-lab-2",
		toEventId: "evt-annual-2026",
		relationshipType: "supports",
		createdAt: "2026-01-15T10:30:00Z",
	},
];

// =============================================================================
// COMBINED GRAPH DATA
// =============================================================================

export const MOCK_GRAPH_DATA: GraphData = {
	events: MOCK_EVENTS,
	edges: MOCK_EDGES,
};

/**
 * Get mock graph data (simulates API call).
 * In production, this would be replaced by actual API fetch.
 */
export function getMockGraphData(): Promise<GraphData> {
	return new Promise((resolve) => {
		setTimeout(() => {
			resolve(MOCK_GRAPH_DATA);
		}, 500); // Simulate network delay
	});
}

/**
 * Get a single event by ID.
 */
export function getMockEventById(id: string): TimelineEvent | undefined {
	return MOCK_EVENTS.find((event) => event.id === id);
}

/**
 * Get edges connected to a specific event.
 */
export function getMockEdgesForEvent(eventId: string): EventEdge[] {
	return MOCK_EDGES.filter(
		(edge) => edge.fromEventId === eventId || edge.toEventId === eventId,
	);
}
