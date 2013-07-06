election
========

A start at implementing an online minimum trust secret ballot election system.

Status
------

In progress, not usable.

Currently, there is a working server, and a working client that implement the full voting protocol,
including all the signatures and almost all of the signature validation.

There are currently no tools for setting up elections, creating lists of voters, ending elections,
sharing the election results etc. Thus the this system is not yet usable.

Executables:
- createVoter.go creates some test voters
- createElection.go creates a test election for the test voters
- server.go runs all the servers 
- client.go casts a ballot as one of the test voters

Background
==========

An `Election`, as considered herein, consists of collecting `ballot`s from a subset of `voters`, and publishing the results.

Basically:
- all `voters` can (but are not required) to `vote` meaning contribute exactly one `ballot` to the `election`.
- the set of `ballot`s contributed to the election publicly visible upon the end of the `election`

There are some additional goals:
- resistance to coercion: it is impossible to force a voter to vote a particular way
- verifiability:
    - (1) voters can guarantee that their `ballot` is included and unedited
    - (2) voters can guarantee that each voter contributes at most one `ballot`
    - (3) voters can guarantee no `ballot`s are included without being contributed by a voter
    - (4) voters can guarantee the set of voters is correct (no one denied, no one wrongly allowed to vote)
    
Unfortunately, in any realistic situation this is impossible. Perfect resistance to coercion requires
voters to not be able to be controlled in any way,
and verifiability requires at minimum trust that the list of potential voters ia correct.

Despite these issues, elections are conducted, and in some cases, are done well enough that the public accepts the results.
My goal for this project is to do at better than existing systems
(particularly paper ballot, and existing electronic systems in the United States), while also lowering the cost.


Particularly, I'll target beating the vote by mail systems, since they seem to be well accepted,
are easy to analyze, and really easy to beat.
This means that this system aims to better resistance to coercion and higher verifiability
than existing vote by mail systems.


Coercion issues
---------------

The first defense against coercion is the secret ballot: the property that one can not associate a ballot with its voter.
Without this, its easy to reward or punish people for voting a particular way. This can fail for trivially small numbers of voters,
but generally seems to be accepted as effective by the public.

The second defense generally seems on the legal side: in some countries, you will get in very serious trouble if you are found
trying to coerce votes.

Its important to note that one should not be able to share their ballot in a manner that can prove that they really voted or will vote that way,
since then people could simply be coerced to do that. This is what voting booths are all about, and this is the main place vote by mail fails.
I'll mention here that my system has about the same issues as vote my mail here (perhaps worse),
but could also be done in a voting booth type setup which would fix the issue. More about that later.

This is a particularly hard problem to solve in a robust manner since it includes vulnerabilities to social engineering.
Threats, bribes, replacing people with replicants, mind control etc. are all attack vectors here. Its not a solvable problem,
but hopefully we can prevent it from being a significant issue in actual elections. When analyzing an electoral system, always
keep coercion in mind, since it will potentially be an issue, you need to evaluate how bad of an issue it is!

Verifiability issues
--------------------

This is where the existing system fail really badly. Of the three factors listed above, existing systems generally
don't even to a partial job of any of them. Its purely done based on trusting the administration of the election.
This is very bad, especially given the elections are often not run by trusted neutral parties.

Of the 4 factors listed above, my system solves (1), and does at least as good at the others.

There are also some nasty issues regarding enforcing when the polls should close. I'm pretty sure this isn't a perfect solution to that,
but its pretty easy to come up with something decent. I we have to ask voters to accept that corrupt management of the ballot server could
introduce some fraud if people are still trying to vote right up to and/or after the deadline regarding which of these votes are rejected.

The Design
==========

- A document (the `election description`) listing the public keys for the various involved servers, and any other needed meta dada (like what is being voted on)
is published
- A set of public keys for registered voters is published. This does not need to (and likely should not be) anonymized.
- Voters cast their votes (process below)
- Election ends
- Votes are displayed (anonymously) in a public table
- Subset of voters that has ballots signed is displayed in a public table (including their signature requests signed with their private keys)

The vote submission process
- Voter constructs a ballot containing their choices
- A random integer is included in the ballot (The only reason here is that if the voter fails to produce a unique ballot, their vote is not counted)
- The ballot is blinded (for RSA blind signatures)
- The blinded ballot is signed with the voter's private key
- The signed blinded ballot is send as a Signature Request to the Ballot server
- The ballot server must respond with Signature Response containing a valid request signed with the same key,
as well as a signature for the blinded ballot. If the voter has gotten a ballot signed previously, their old request will be included
(which proves they have voted previously while also providing a copy of the signed ballot incase it was lost)
- The voter unblinds their ballot's signature and anonymously submits it and the ballot to the Vote server
- The vote server must respond with a Ballot Entry, which contains the ballot, and is signed with the Vote key.

When the final votes are listed, the Voter may check that their vote is included (since it is unique, they can find it in the list).
If their vote is not included, they can provide the signed vote they have from the vote server as proof their vote should have been included (and thus invalidating the election).

In the number of votes in the final votes table is larger than the number of voters from which the ballot server can provide valid signature requests (signed with their private keys),
the election is shown to be invalid.

If the required response to a request from a voter is not given by one of the servers, the voter can proxy the request through any third party (such as news media or auditors) which will either
get them a valid response or provide evidence that the election is invalid. Since all the requests are idempotent, and anonymity  and security are not lost by third parties viewing the messages, this is safe.

The only major concern as far as secret ballot anonymity is determining the origin of the final vote submission. To solve this, submission through Tor is recommended. There is still potential for timing attacks,
so some variable delay between getting the ballot signed and submitted is also recommended. The election could be broken up (temporally) into 2 stages to solve this.