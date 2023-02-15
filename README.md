<div align="center"><img src="./docs/images/charonlogo.svg" /></div>
<h1 align="center">Charon<br/>განაწილებადი ვალიდატორს შუალედური პროგრამული უზრუნველყოფის კლიენტი</h1>

<p align="center"><a href="https://github.com/obolnetwork/charon/releases/"><img src="https://img.shields.io/github/tag/obolnetwork/charon.svg"></a>
<a href="https://github.com/ObolNetwork/charon/blob/main/LICENSE"><img src="https://img.shields.io/github/license/obolnetwork/charon.svg"></a>
<a href="https://godoc.org/github.com/obolnetwork/charon"><img src="https://godoc.org/github.com/obolnetwork/charon?status.svg"></a>
<a href="https://goreportcard.com/report/github.com/obolnetwork/charon"><img src="https://goreportcard.com/badge/github.com/obolnetwork/charon"></a>
<a href="https://github.com/ObolNetwork/charon/actions/workflows/golangci-lint.yml"><img src="https://github.com/obolnetwork/charon/workflows/golangci-lint/badge.svg"></a></p>

ეს საცავი შეიცავს Charon ვალიდიატორის განაწილებულადი კლიენტის წყაროს კოდს (გამოითქმის "kharon"); HTTP შუა პროგრამის კლიენტი Ethereum Staking-ისთვის, რომელიც საშუალებას გაძლევთ უსაფრთხოდ გაუშვათ ერთი ვალიდატორი დამოუკიდებელი კვანძების ჯგუფზე.

Charon-ს ახლავს ვებ აპლიკაცია სახელწოდებით [Distributed Validator Launchpad](https://goerli.launchpad.obol.tech/) განაწილებულადი ვალიდატორის გასაღების გენერირებისთვის.

სტეიკერები იყენებენ charon-ს Ethereum ვალიდატორების მუშაობის პასუხისმგებლობის გასანაწილებლად რამდენიმე სხვადასხვა ინსტანციაზე და კლიენტის იმპლემენტაციაზე.

![Example Obol Cluster](./docs/images/DVCluster.png)

###### განაწილებულადი ვალიდატორის კლასტერი, რომელიც იყენებს Charon კლიენტს კლიენტისა და ტექნიკის უკმარისობის რისკების დასაცავად.

## სწრაფი დაწყება

charon-ის შესამოწმებლად უმარტივესი გზაა[charon-distributed-validator-cluster](https://github.com/ObolNetwork/charon-distributed-validator-cluster) რეპო
რომელიც შეიცავს დოკერის შედგენის დაყენებას, რათა აწარმოოს სრული ქარონის კლასტერი თქვენს ადგილობრივ აპარატზე.

## დოკუმენტაცია

[Obol Docs](https://docs.obol.tech/) ვებგვერდი საუკეთესო საშუალებაა გასაწყებად.
მნიშვნელოვანი სექცია არის [intro](https://docs.obol.tech/docs/intro),
[key concepts](https://docs.obol.tech/docs/int/key-concepts) და [charon](https://docs.obol.tech/docs/dv/introducing-charon).

დეტალური დოკუმენტაციისათვის ამ რეპოში,იხილეთ [docs](docs) საქაღალდე:

- [Configuration](docs/configuration.md): charon node კონფიგურაცია
- [Architecture](docs/architecture.md): charon cluster-ის მიმოხილვა და node-ის არქიტურა
- [Project Structure](docs/structure.md): პროექტის საქაღალდის სტრუქტურა
- [Branching and Release Model](docs/branching.md): Git ფილიალი და გამოშვების მოდელი
- [Go Guidelines](docs/goguidelines.md): გაიდლაინები და პრინციპები, რომლებიც დაკავშირებულია GO-ს განვითარებასთან
- [Contributing](docs/contributing.md): როგორ შევიტანოთ წვლილი charon-ში; githuks, PR შაბლონები და ა.შ.

ყოველთვის არის [charon godocs](https://pkg.go.dev/github.com/obolnetwork/charon) სარესურსო კოდების დოკუმენტაციისათვის.

## Supported Consensus Layer Clients

Charon integrates into the Ethereum consensus stack as a middleware between the validator client
and the beacon node via the official [Eth Beacon Node REST API](https://ethereum.github.io/beacon-APIs/#/).
Charon supports any upstream beacon node that serves the Beacon API.
Charon aims to support any downstream standalone validator client that consumes the Beacon API.

| Client                                             | Beacon Node | Validator Client | Notes                                   |
| -------------------------------------------------- | :---------: | :--------------: |-----------------------------------------|
| [Teku](https://github.com/ConsenSys/teku)          |     ✅      |        ✅        | Fully supported                         |
| [Lighthouse](https://github.com/sigp/lighthouse)   |     ✅      |        ✅        | Fully supported                         |
| [Lodestar](https://github.com/ChainSafe/lodestar)  |     ✅      |       \*️⃣        | DVT compatibility issue                 |
| [Vouch](https://github.com/attestantio/vouch)      |     \*️⃣     |        ✅        | Only validator client provided          |
| [Prysm](https://github.com/prysmaticlabs/prysm)    |     ✅      |        🛑        | Validator client requires gRPC API      |
| [Nimbus](https://github.com/status-im/nimbus-eth2) |     ✅      |        ✅        | Soon to be supported |

## Project Status

It is still early days for the Obol Network and things are under active development.
We are moving fast so check back in regularly to track the progress.

Charon is a distributed validator, so its main responsibility is performing validation duties.
The following table outlines which clients have produced which duties on a public testnet, and which are still under construction (🚧 )

| Duty \ Client                        |                      Teku                      |                    Lighthouse                    | Lodestar | Nimbus | Vouch | Prysm |
|--------------------------------------|:----------------------------------------------:|:------------------------------------------------:|:--------:|:------:|:-----:|:-----:|
| _Attestation_                        |                       ✅                        |                        ✅                         |    🚧    |   🚧   |  ✅   |  🚧   |
| _Attestation Aggregation_            |                       🚧                       |                        🚧                        |    🚧    |   🚧   |  🚧   |  🚧   |
| _Block Proposal_                     |                       ✅                        |                        ✅                         |    🚧    |   🚧   |  🚧   |  🚧   |
| _Blinded Block Proposal (mev-boost)_ | [✅](https://ropsten.beaconcha.in/block/555067) | [✅](https://ropsten.etherscan.io/block/12822070) |    🚧    |   🚧   |  🚧   |  🚧   |
| _Sync Committee Message_             |                       ✅                        |                        ✅                         |    🚧    |   🚧   |  🚧   |  🚧   |
| _Sync Committee Contribution_        |                       🚧                       |                        🚧                        |    🚧    |   🚧   |  🚧   |  🚧   |
