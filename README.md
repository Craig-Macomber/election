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
- A set of public keys for registered voters is collected. This does not need to (and likely should not be) anonymized.
- A document (the `election description`) listing the public keys for the various involved servers, all the voter public keys, and any other needed meta dada (like what is being voted on)
is published. If there are many voters, a hash of the list of voters (or the root of a tree of hashes) can be included instead of the full list (which would be published separately)
- Voters cast their votes (process below)
- Election ends
- Votes are displayed (anonymously) in a public table
- Subset of voters that has ballots signed is displayed in a public table (including a signed subsection of their signature requests signed with their private keys. This does not include the blinded ballot!)

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

Diagrams
--------
When interacting with one of the servers, this process is followed:
![Voting Process](/documentation/ElectionServerRequest_Flow.png "Voting Process")

Since all operations with the servers are idempotent, and none (in either the request or responses) expose both the content of your vote and your identity, its perfectly safe to let others try and submit your requests. If they are malicious, it does no harm. This approach can be used to expose a misbeahving election server: if your ballot signing request is not processed correctly for example, you could have someone else, say the UN, submit it for you, and forward the response. If your request is still not handeled correctly, they now have proof your request has wrongfully been denied, and otherwise you have our needed response and can continue.



The entire process of casting a vote in the election from a voter's perspective:
![Voting Process](/documentation/Election_Flow.png "Voting Process")


Attacks
=======

Here is a brief summery of some of the known attacks, including ones that have been mitigated.

Keep in mind that the goal of this project is to implement a system better than main in ballots, which also suffer from many very serious issues.

Currently all attacks here are related to coercion, which is an unsolvable problem.
No flaws in the security allowing stealing a vote, or adding votes or removing votes have been found so far.

Mitigated Attacks
-----------------
Attacks that have been found and fixed so far.

Blinding factor disclosure coercion attack: In the original design,
the blinded ballots were displayed as part of the signature request for each voter after the election.
This enabled an attacker to be able to request that a voter use software that allowed extracting blinding factor,
so the voter could be forced/asked to disclose this, which would be sufficient to prove which vote was theirs.
This allows the attacker to punish or reward voters based on their vote, since proof of how they voted can be provided.
The fix is to not display the blinded ballots.

This attack is strictly a coercion attack: it helps people to illegally coercion voters, but it does not enable any other kind of fraud.
As discussed above, preventing coercion attacks completely is impossible, but this one was particularly easy to exploit in a wide spread manner
(robust pay per vote would be easy) and easy to fix.

Unsolved attacks
----------------
Private key selling/extortion: A voter can be forced to disclose their private key, or can sell it
(and thus their vote assuming they don't vote first, which their buyer could complain about).

This is the equivalent of selling your ballot in a mail in system. In mail in systems (such as in Washington State) where you must physically sign the envelope,
this is equivalent you selling your ballot along with a pre-signed envelope.

It appears that a voter needs to be able to verify that his public key is included in the election, which apparently
necessitates that an attacker can, if given a public key, can verify if its valid or not (at the very least, they can try and vote with it).
Thus it appears that there is no possible complete mitigation for this attack.
Some delayed response/verification from the servers could help reduce the ease/effectiveness of the attack (but not by a lot), but that complicates validation and hurts ease of use drastically.


Paid ballot signing Attack:
The voter provides a payment address (say bitcoin) to the attacker and requests a ballot from the attacker.
The attacker provides the ballot. The voter gets it signed (consuming their right to vote with it).
The voter then provides this ballot to the attacker to prove they got it signed. The attacker (and optionally the voter) cast the ballot.
Once the attacker has proof the ballot will be included (with this election system, this is as soon as they get the signed ballot, but they could wait until after the elections)
the attacker pays the voter.

This attack, with proper use of bitcoin and Tor would be easy to do in a manner that is automated for the attacker, anonymous for the voter, and pseudonymous for the attacker (via Tor hidden service).
This means the attacker could get a reputation of actually paying out over multiple elections and thus could effectively buy votes with no danger of getting caught for them or the voters.

This is a very serious and crippling attack. Currently a solution is not known. If you have suggestions, let me know.

The same attack could be done via private key selling, but this would endanger the voter (the attacker would know who they were and could report it).
 
