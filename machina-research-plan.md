# Machina Labs Research Plan

**Purpose:** Prepare for in-person visit on April 3, 2026. Lunch + tour with Jeremy Boyle at the main Machina office, then visit to second facility to meet Gary Luo (software engineer, more relevant to evaluating Miru as a tool for Machina).

---

## What We Already Know

### About Machina Labs
- **Founded:** 2019, HQ in Chatsworth, CA
- **CEO/Co-Founder:** Edward Mehr (ex-SpaceX, Relativity, USC)
- **Co-Founder:** Babak Raeisinia (Head of Applications & Partnerships)
- **VP Engineering:** Kyle Hickey
- **VP Production:** John Borrego
- **Employees:** ~89, growing fast after Series C
- **Total funding:** ~$174M, including $124M Series C (Feb 2026) led by Woven Capital (Toyota), Lockheed Martin Ventures, Balerion Space Ventures, Strategic Development Fund
- **Core tech:** RoboCraftsman platform — AI-driven robotic sheet metal forming (no dies/molds). Turns CAD into precision-formed metal parts in hours instead of months.
- **Industries:** Defense/aerospace (USAF contracts, defense primes for missiles/hypersonics), automotive (Toyota partnership), electronics, heavy machinery
- **Facilities:**
  - Main office: 9410 Owensmouth Ave, Chatsworth, CA
  - Second facility: 20559 Prairie St, Chatsworth (~60,300 sq ft manufacturing warehouse, added 2024)
  - Planned: 200,000 sq ft "Intelligent Factory" (~50 RoboCraftsman robots)

### About Our Contacts
- **Jeremy Boyle** — Mechanical Engineer. Masters from UT Austin (your TIW connection). Aviation Week 20 Twenties 2024 winner. Previously interned at Apple (Manufacturing Design). Now at Machina working on metal forming robotics.
- **Gary Luo (Weijia Luo)** — Software Engineer. BS in CS from UIUC (2011-2016). Previously at Yahoo, Enigma Technologies, Bikky. GitHub: garyluoex. More relevant contact for evaluating Miru adoption — likely involved in software infrastructure decisions for their robotic fleet.

---

## Research Plan

### 1. Understand Machina's Software/Fleet Architecture
**Why:** This is the key to determining if Miru is a fit. Machina operates fleets of robots (RoboCraftsmen) that need to be configured per-job, per-material, per-part-geometry.

**Research tasks:**
- [ ] Search for Machina Labs engineering blog posts or technical talks about their software stack
- [ ] Look for conference presentations (ROSCon, ICRA, manufacturing conferences) by Machina engineers
- [ ] Search GitHub for any open-source repos under Machina-Labs org
- [ ] Check if they use ROS/ROS2 or a custom robotics framework
- [ ] Research what "software-defined factory" means in their context — how is configuration managed today?

### 2. Map the Stakeholder Landscape
**Why:** Understand who influences technology adoption decisions beyond Gary.

**Research tasks:**
- [ ] Deep-dive on Kyle Hickey (VP Engineering) — background, prior companies, what he cares about technically
- [ ] Research Babak Raeisinia (Co-Founder, Head of Apps & Partnerships) — he may own partner/vendor decisions
- [ ] Research John Borrego (VP Production) — production scaling needs drive config management pain
- [ ] Look for other software/infrastructure engineers on LinkedIn or The Org chart
- [ ] Identify if they have a DevOps/infra team or if software engineers own deployment
- [ ] Check for any hiring posts that reveal their tech stack or pain points (e.g., "looking for someone to manage robot fleet configs")

### 3. Understand Their Scaling Pain Points
**Why:** Machina is about to 3x+ their robot fleet with the new Intelligent Factory. This is exactly when config management becomes critical.

**Research tasks:**
- [ ] Read the Series C press coverage in detail for hints about scaling challenges
- [ ] Research what config/parameter management looks like for industrial robotic workcells
- [ ] Look into how similar companies (Relativity Space, Bright Machines, Rapid Robotics) handle fleet config
- [ ] Understand what parameters a sheet-forming robot needs per job (force profiles, tool paths, material properties, sensor calibration, etc.)

### 4. Understand Their Customer/Contract Requirements
**Why:** Defense and aerospace customers (USAF, Lockheed) have strict traceability and audit requirements — a natural fit for Miru's audit trail features.

**Research tasks:**
- [ ] Research AS9100, ITAR, and CMMC compliance requirements for defense manufacturing software
- [ ] Look into what traceability/audit requirements exist for defense parts manufacturing
- [ ] Check if Machina has any published case studies or press about their defense work
- [ ] Research NCMS (National Center for Manufacturing Sciences) member spotlight for details on their defense programs

### 5. Assess the Miru Value Prop Fit
**Why:** Prepare specific talking points for the Gary Luo meeting.

**Research tasks:**
- [ ] Map Miru capabilities to Machina pain points:
  - Config schemas -> Robot job parameter validation
  - Audit trail -> Defense traceability requirements
  - Fleet management -> Managing 50+ RoboCraftsmen across facilities
  - Release management -> Software updates across robot fleet
  - Validation before deployment -> Catching bad configs before they ruin a $10K sheet of titanium
- [ ] Identify what Miru does NOT currently support that Machina might need (real-time parameter tuning? closed-loop config updates from sensor data?)
- [ ] Prepare questions to ask Gary about their current workflow for deploying configs to robots

### 6. Prepare Conversation Topics for Jeremy
**Why:** Jeremy is the warm intro and relationship. Keep this meeting collegial and informative.

**Research tasks:**
- [ ] Review his Aviation Week 20 Twenties profile for talking points
- [ ] Understand what a Mechanical Engineer's day-to-day looks like at Machina (likely process development, robot cell design)
- [ ] Prepare questions about what it's like working there, growth trajectory, team culture
- [ ] Think about UT Austin shared context / TIW connections

---

## Prioritized Question List to Develop

### For Jeremy (lunch/tour)
1. How has the team grown since you joined? What's the culture like?
2. How many RoboCraftsman cells are running today?
3. What does a typical job flow look like end-to-end (CAD in -> part out)?
4. What's the most challenging part of scaling production?
5. How do different teams (software, ME, process) collaborate on a new part program?

### For Gary (second facility)
1. What does your software stack look like? (ROS? Custom? Cloud?)
2. How do you currently manage configuration across your robot fleet?
3. When you add a new RoboCraftsman cell, what's the setup/config process?
4. What's the biggest software pain point as you scale to the Intelligent Factory?
5. How do you handle config versioning and rollback if something goes wrong?
6. What are your traceability requirements from defense/aero customers?
7. Have you looked at any config management tools, or is it all custom/in-house?

---

## Timeline

| When | Action |
|------|--------|
| **Today (April 2)** | Execute research tasks above |
| **Tonight** | Synthesize findings into a 1-pager of talking points |
| **Tomorrow AM** | Review talking points, prep questions |
| **Tomorrow lunch** | Meet Jeremy, tour main facility, learn organically |
| **Tomorrow afternoon** | Meet Gary at second facility, deeper technical conversation |
| **Post-visit** | Follow up with tailored Miru pitch based on what we learned |
