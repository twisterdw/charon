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

## მხარდაჭერილი კონსენსუსის ფენის კლიენტები

Charon ინტეგრირდება Ethereum-ის კონსენსუსის სტეკში, როგორც შუალედური პროგრამა ვალიდატორ კლიენტს შორის
და შუქურა კვანძის მეშვეობით ოფიციალური [Eth Beacon Node REST API](https://ethereum.github.io/beacon-APIs/#/).
ქარონი მხარს უჭერს ნებისმიერ ზედა დინების შუქურას, რომელიც ემსახურება Beacon API-ს.
Charon მიზნად ისახავს მხარი დაუჭიროს ნებისმიერ ქვემოთ მოყვანილ დამოუკიდებელ ვალიდატორ კლიენტს Beacon API-ის გამოყენებით.

| Client                                             | Beacon Node | Validator Client | Notes                                   |
| -------------------------------------------------- | :---------: | :--------------: |-----------------------------------------|
| [Teku](https://github.com/ConsenSys/teku)          |     ✅      |        ✅        | სრულად მხარდაჭერი                     |
| [Lighthouse](https://github.com/sigp/lighthouse)   |     ✅      |        ✅        | სრულად მხარდაჭერილი                  |
| [Lodestar](https://github.com/ChainSafe/lodestar)  |     ✅      |       \*️⃣        | თავსებადობის პრობლემა DVT-თან          |
| [Vouch](https://github.com/attestantio/vouch)      |     \*️⃣     |        ✅        | მოწოდებულია მხოლოდ ვალიდატორი კლიენტი|
| [Prysm](https://github.com/prysmaticlabs/prysm)    |     ✅      |        🛑        | Validator კლიენტი საჭიროებს gRPC API-ს |
| [Nimbus](https://github.com/status-im/nimbus-eth2) |     ✅      |        ✅        | მალე იქნება მხარდაჭერილი |

## პროექტის სტატუსი

Obol Network-ისთვის ჯერ ადრეა და ყველაფერი აქტიური განვითარების პროცესშია.
ჩვენ სწრაფად მივიწევთ წინ, ამიტომ რეგულარულად შეამოწმეთ პროგრესი.

ქარონი არის განაწილებადი ვალიდატორი, ამიტომ მისი მთავარი პასუხისმგებლობაა ვალიდაციის მოვალეობების შესრულება.
შემდეგი ცხრილი ასახავს კლიენტებს, რომელი მოვალეობები შექმნეს საჯარო საცდელ ქსელზე და რომელი ჯერ კიდევ მშენებლობის პროცესშია (🚧 )

| Duty \ Client                        |                      Teku                      |                    Lighthouse                    | Lodestar | Nimbus | Vouch | Prysm |
|--------------------------------------|:----------------------------------------------:|:------------------------------------------------:|:--------:|:------:|:-----:|:-----:|
| _Attestation_                        |                       ✅                        |                        ✅                         |    🚧    |   🚧   |  ✅   |  🚧   |
| _Attestation Aggregation_            |                       🚧                       |                        🚧                        |    🚧    |   🚧   |  🚧   |  🚧   |
| _Block Proposal_                     |                       ✅                        |                        ✅                         |    🚧    |   🚧   |  🚧   |  🚧   |
| _Blinded Block Proposal (mev-boost)_ | [✅](https://ropsten.beaconcha.in/block/555067) | [✅](https://ropsten.etherscan.io/block/12822070) |    🚧    |   🚧   |  🚧   |  🚧   |
| _Sync Committee Message_             |                       ✅                        |                        ✅                         |    🚧    |   🚧   |  🚧   |  🚧   |
| _Sync Committee Contribution_        |                       🚧                       |                        🚧                        |    🚧    |   🚧   |  🚧   |  🚧   |
