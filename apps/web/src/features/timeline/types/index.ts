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
		bg: "var(--event-consultation-bg)",
		border: "var(--event-consultation-border)",
		text: "var(--event-consultation-text)",
		glow: "var(--event-consultation-glow)",
	},
	diagnosis: {
		bg: "var(--event-diagnosis-bg)",
		border: "var(--event-diagnosis-border)",
		text: "var(--event-diagnosis-text)",
		glow: "var(--event-diagnosis-glow)",
	},
	prescription: {
		bg: "var(--event-prescription-bg)",
		border: "var(--event-prescription-border)",
		text: "var(--event-prescription-text)",
		glow: "var(--event-prescription-glow)",
	},
	procedure: {
		bg: "var(--event-procedure-bg)",
		border: "var(--event-procedure-border)",
		text: "var(--event-procedure-text)",
		glow: "var(--event-procedure-glow)",
	},
	lab_result: {
		bg: "var(--event-lab-bg)",
		border: "var(--event-lab-border)",
		text: "var(--event-lab-text)",
		glow: "var(--event-lab-glow)",
	},
	imaging: {
		bg: "var(--event-imaging-bg)",
		border: "var(--event-imaging-border)",
		text: "var(--event-imaging-text)",
		glow: "var(--event-imaging-glow)",
	},
	note: {
		bg: "var(--event-note-bg)",
		border: "var(--event-note-border)",
		text: "var(--event-note-text)",
		glow: "var(--event-note-glow)",
	},
	vaccination: {
		bg: "var(--event-vaccination-bg)",
		border: "var(--event-vaccination-border)",
		text: "var(--event-vaccination-text)",
		glow: "var(--event-vaccination-glow)",
	},
	allergy: {
		bg: "var(--event-allergy-bg)",
		border: "var(--event-allergy-border)",
		text: "var(--event-allergy-text)",
		glow: "var(--event-allergy-glow)",
	},
	visit_note: {
		bg: "var(--event-note-bg)",
		border: "var(--event-note-border)",
		text: "var(--event-note-text)",
		glow: "var(--event-note-glow)",
	},
	vital_signs: {
		bg: "var(--event-vital-bg)",
		border: "var(--event-vital-border)",
		text: "var(--event-vital-text)",
		glow: "var(--event-vital-glow)",
	},
	referral: {
		bg: "var(--event-referral-bg)",
		border: "var(--event-referral-border)",
		text: "var(--event-referral-text)",
		glow: "var(--event-referral-glow)",
	},
	insurance_claim: {
		bg: "var(--event-insurance-bg)",
		border: "var(--event-insurance-border)",
		text: "var(--event-insurance-text)",
		glow: "var(--event-insurance-glow)",
	},
	other: {
		bg: "var(--event-other-bg)",
		border: "var(--event-other-border)",
		text: "var(--event-other-text)",
		glow: "var(--event-other-glow)",
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
