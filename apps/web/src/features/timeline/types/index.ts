/**
 * Timeline Types
 *
 * Protocol-aligned types for medical timeline events and graph relationships.
 * These types mirror the backend entities in apps/backend/internal/timeline/entity.go
 */

// =============================================================================
// EVENT TYPES
// =============================================================================

/**
 * Supported medical event types in the timeline.
 * Matches backend: TimelineEventType in entity.go
 */
export const TimelineEventType = {
	CONSULTATION: "consultation",
	DIAGNOSIS: "diagnosis",
	PRESCRIPTION: "prescription",
	PROCEDURE: "procedure",
	LAB_RESULT: "lab_result",
	IMAGING: "imaging",
	NOTE: "note",
	VACCINATION: "vaccination",
	ALLERGY: "allergy",
	VISIT_NOTE: "visit_note",
	VITAL_SIGNS: "vital_signs",
	REFERRAL: "referral",
	INSURANCE_CLAIM: "insurance_claim",
	OTHER: "other",
} as const;

export type TimelineEventType =
	(typeof TimelineEventType)[keyof typeof TimelineEventType];

/**
 * Human-readable labels for each event type.
 */
export const EVENT_TYPE_LABELS: Record<TimelineEventType, string> = {
	consultation: "Consultation",
	diagnosis: "Diagnosis",
	prescription: "Prescription",
	procedure: "Procedure",
	lab_result: "Lab Results",
	imaging: "Imaging",
	note: "Note",
	vaccination: "Vaccination",
	allergy: "Allergy",
	visit_note: "Visit Note",
	vital_signs: "Vital Signs",
	referral: "Referral",
	insurance_claim: "Insurance Claim",
	other: "Other",
};

/**
 * Lucide icon names associated with each event type.
 */
export const EVENT_TYPE_ICONS: Record<TimelineEventType, string> = {
	consultation: "Stethoscope",
	diagnosis: "ClipboardList",
	prescription: "Pill",
	procedure: "Scissors",
	lab_result: "TestTube2",
	imaging: "ScanLine",
	note: "FileText",
	vaccination: "Syringe",
	allergy: "AlertTriangle",
	visit_note: "FileText",
	vital_signs: "HeartPulse",
	referral: "UserPlus",
	insurance_claim: "FileCheck",
	other: "File",
};

/**
 * Color themes for each event type (for graph nodes)
 * Dark mode optimized: transparent tinted backgrounds, bright borders, light text
 */
export const EVENT_TYPE_COLORS: Record<
	TimelineEventType,
	{ bg: string; border: string; text: string; glow: string }
> = {
	consultation: {
		bg: "rgba(59, 130, 246, 0.15)",
		border: "#3b82f6",
		text: "#dbeafe",
		glow: "rgba(59, 130, 246, 0.3)",
	},
	diagnosis: {
		bg: "rgba(239, 68, 68, 0.15)",
		border: "#ef4444",
		text: "#fee2e2",
		glow: "rgba(239, 68, 68, 0.3)",
	},
	prescription: {
		bg: "rgba(34, 197, 94, 0.15)",
		border: "#22c55e",
		text: "#dcfce7",
		glow: "rgba(34, 197, 94, 0.3)",
	},
	procedure: {
		bg: "rgba(245, 158, 11, 0.15)",
		border: "#f59e0b",
		text: "#fef3c7",
		glow: "rgba(245, 158, 11, 0.3)",
	},
	lab_result: {
		bg: "rgba(99, 102, 241, 0.15)",
		border: "#6366f1",
		text: "#e0e7ff",
		glow: "rgba(99, 102, 241, 0.3)",
	},
	imaging: {
		bg: "rgba(168, 85, 247, 0.15)",
		border: "#a855f7",
		text: "#f3e8ff",
		glow: "rgba(168, 85, 247, 0.3)",
	},
	note: {
		bg: "rgba(100, 116, 139, 0.15)",
		border: "#64748b",
		text: "#f1f5f9",
		glow: "rgba(100, 116, 139, 0.3)",
	},
	vaccination: {
		bg: "rgba(6, 182, 212, 0.15)",
		border: "#06b6d4",
		text: "#cffafe",
		glow: "rgba(6, 182, 212, 0.3)",
	},
	allergy: {
		bg: "rgba(234, 179, 8, 0.15)",
		border: "#eab308",
		text: "#fef9c3",
		glow: "rgba(234, 179, 8, 0.3)",
	},
	visit_note: {
		bg: "rgba(100, 116, 139, 0.15)",
		border: "#64748b",
		text: "#f1f5f9",
		glow: "rgba(100, 116, 139, 0.3)",
	},
	vital_signs: {
		bg: "rgba(236, 72, 153, 0.15)",
		border: "#ec4899",
		text: "#fce7f3",
		glow: "rgba(236, 72, 153, 0.3)",
	},
	referral: {
		bg: "rgba(139, 92, 246, 0.15)",
		border: "#8b5cf6",
		text: "#ede9fe",
		glow: "rgba(139, 92, 246, 0.3)",
	},
	insurance_claim: {
		bg: "rgba(16, 185, 129, 0.15)",
		border: "#10b981",
		text: "#ecfdf5",
		glow: "rgba(16, 185, 129, 0.3)",
	},
	other: {
		bg: "rgba(115, 115, 115, 0.15)",
		border: "#737373",
		text: "#f5f5f5",
		glow: "rgba(115, 115, 115, 0.3)",
	},
};

// =============================================================================
// RELATIONSHIP TYPES
// =============================================================================

/**
 * Defines how two events are connected in the graph.
 * Matches backend: RelationshipType in entity.go
 */
export const RelationshipType = {
	RESULTED_IN: "resulted_in", // consultation → diagnosis
	LEAD_TO: "lead_to", // diagnosis → treatment
	REQUESTED_BY: "requested_by", // lab test ← consultation
	SUPPORTS: "supports", // lab report → diagnosis
	FOLLOWS_UP: "follows_up", // visit → previous visit
	CONTRADICTS: "contradicts", // finding → previous finding
	ATTACHED_TO: "attached_to", // file → event
} as const;

export type RelationshipType =
	(typeof RelationshipType)[keyof typeof RelationshipType];

/**
 * Human-readable labels for relationship types.
 */
export const RELATIONSHIP_LABELS: Record<RelationshipType, string> = {
	resulted_in: "Resulted In",
	lead_to: "Led To",
	requested_by: "Requested By",
	supports: "Supports",
	follows_up: "Follows Up",
	contradicts: "Contradicts",
	attached_to: "Attached To",
};

// =============================================================================
// ENTITIES
// =============================================================================

/**
 * Represents a single event in the patient's medical timeline.
 * This is a node in the timeline graph.
 */
export interface TimelineEvent {
	readonly id: string;
	readonly patientId: string;
	readonly type: TimelineEventType;
	readonly title: string;
	readonly description?: string;
	readonly provider?: string;
	readonly timestamp: string; // ISO string from API
	readonly blobRef?: string;
	readonly isEncrypted: boolean;
	readonly metadata?: Record<string, unknown>;
	readonly createdAt: string;
	readonly updatedAt: string;
}

/**
 * Represents a directed relationship between two timeline events.
 * This is an edge in the timeline graph.
 */
export interface EventEdge {
	readonly id: string;
	readonly fromEventId: string;
	readonly toEventId: string;
	readonly relationshipType: RelationshipType;
	readonly metadata?: Record<string, unknown>;
	readonly createdAt: string;
}

/**
 * Represents a file attached to an event.
 */
export interface EventFile {
	readonly id: string;
	readonly eventId: string;
	readonly blobRef: string;
	readonly fileName: string;
	readonly mimeType: string;
	readonly fileSize: number;
	readonly metadata?: Record<string, unknown>;
	readonly createdAt: string;
}

// =============================================================================
// API RESPONSES
// =============================================================================

/**
 * Response from /timeline/graph endpoint.
 * Contains all events and edges for visualization.
 */
export interface GraphData {
	readonly events: TimelineEvent[];
	readonly edges: EventEdge[];
}

/**
 * Response from event creation.
 */
export interface UploadResponse {
	readonly success: boolean;
	readonly event?: TimelineEvent;
	readonly message?: string;
}

/**
 * Response from linking events.
 */
export interface LinkResponse {
	readonly success: boolean;
	readonly edge?: EventEdge;
	readonly message?: string;
}

/**
 * Response from related events query.
 */
export interface RelatedEventsResponse {
	readonly events: TimelineEvent[];
}
