package box

// CaptainsDecisionRule is the clause every app line that hands the Mate a
// Crew's question carries (query.BoxResolveLine, autopilot.Line): an open
// decision the brief gives the captain is not the Mate's to answer. On mihu
// 2026-10-14 the `resolve:` line said "decide", and the Mate decided a
// `decides: captain` question on that word although its wake skill said
// otherwise (CON-93): a rule in the tool output is the one a Mate follows.
const CaptainsDecisionRule = "a decides: captain choice goes to the captain as news"
