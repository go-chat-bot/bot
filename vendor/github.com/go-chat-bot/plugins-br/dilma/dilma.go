package dilma

import (
	"fmt"
	"math/rand"
	"regexp"

	"github.com/go-chat-bot/bot"
)

const (
	pattern = "(?i)\\b(dilma)\\b"
)

var (
	re          = regexp.MustCompile(pattern)
	frasesDilma = []string{
		"Primeiro eu queria cumprimentar os internautas. -Oi Internautas! Depois dizer que o meio ambiente é sem dúvida nenhuma uma ameaça ao desenvolvimento sustentável. E isso significa que é uma ameaça pro futuro do nosso planeta e dos nossos países. O desemprego beira 20%, ou seja, 1 em cada 4 portugueses.",
		"No meu xinélo da humildade eu gostaria muito de ver o Neymar e o Ganso. Por que eu acho que.... 11 entre 10 brasileiros gostariam. Você veja, eu já vi, parei de ver. Voltei a ver, e acho que o Neymar e o Ganso têm essa capacidade de fazer a gente olhar.",
		"A única área que eu acho, que vai exigir muita atenção nossa, e aí eu já aventei a hipótese de até criar um ministério. É na área de... Na área... Eu diria assim, como uma espécie de analogia com o que acontece na área agrícola.",
		"Ai você fala o seguinte: \"- Mas vocês acabaram isso?\" Vou te falar: -\"Não, está em andamento!\" Tem obras que \"vai\" durar pra depois de 2010. Agora, por isso, nós já não desenhamos, não começamos a fazer projeto do que nós \"podêmo fazê\"? 11, 12, 13, 14... Por que é que não?",
		"A população ela precisa da Zona Franca de Manau, porque na Zona franca de Manaus, não é uma zona de exportação, é uma zona para o Brasil. Portanto ela tem um objetivo, ela evita o desmatamento, que é altamente lucravito. Derrubar arvores da natureza é muito lucrativo!",
		"Se hoje é o dia das crianças... Ontem eu disse: o dia da criança é o dia da mãe, dos pais, das professoras, mas também é o dia dos animais, sempre que você olha uma criança, há sempre uma figura oculta, que é um cachorro atrás. O que é algo muito importante!",
		"Todos as descrições das pessoas são sobre a humanidade do atendimento, a pessoa pega no pulso, examina, olha com carinho. Então eu acho que vai ter outra coisa, que os médicos cubanos trouxeram pro brasil, um alto grau de humanidade.",
		"Eu dou dinheiro pra minha filha. Eu dou dinheiro pra ela viajar, então é... é... Já vivi muito sem dinheiro, já vivi muito com dinheiro. -Jornalista: Coloca esse dinheiro na poupança que a senhora ganha R$10 mil por mês. -Dilma: O que que é R$10 mil?",
		"Eu já testei e ela [a bola] quica. Eu testei, eu fiz assim uma embaixadinha, minto, uma meia embaixadinha",
		"Aqui, hoje, eu estou saudando a mandioca. Acho uma das maiores conquistas do Brasil",
		"Pra mim essa bola é um símbolo da nossa evolução. Quando nós criamos uma bola dessas, nós nos transformamos em Homo sapiens ou 'mulheres sapiens",
		"Eu acredito que há um pouco de viés sexual ou viés de gênero. Sou descrita como uma mulher dura e forte que coloca o nariz em tudo em que não devia, e eu estou [me dizem] cercada por homens muito fofos",
		"Eu, para ir, eu faço uma escala. Para voltar, eu faço duas, para voltar para o Brasil. Neste caso agora nós tínhamos uma discussão. Eu tinha que sair de Zurique, podia ir para Boston, ou pra Boston, até porque... vocês vão perguntar, mas é mais longe? Não é não, a Terra é curva, viu?",
		"Esse país foi descoberto, foi colonizado através das estradas de água. Essas estradas de água são a forma mais barata de transporte",
		"“A Zona Franca de Manaus, ela está numa região, ela é o centro dela porque é a capital da Amazônia (…). Portanto, ela tem um objetivo, ela evita o desmatamento, que é altamente lucrativo – derrubar árvores plantadas pela natureza é altamente lucrativo",
		"Nós sabemos, e eu já estive aqui várias vezes antes, que essa é uma cidade arborizada, cercada por rios, e que tem essa interessantíssima característica de ter muitas mangueiras. De fato, deve ser muito bom morar numa cidade que de repente você pode chuchar uma árvore e cair uma manga na sua mão. É de fato algo que todo mundo quer, é pegar e ter acesso a uma boa manga",
		"É interessante que muitas vezes no Brasil, você é, como diz o povo brasileiro, muitas vezes você é criticado por ter o cachorro e, outras vezes, por não ter o mesmo cachorro. Esta é uma crítica interessante que acontece no Brasil",
		"O bacalhau é uma moleza de fazer. Posso falar, é simplíssimo o bacalhau. Você corta várias coisas, bota uma camada, bota outra, bota, você vai ver o bacalhau… agora, é sem reclamações, sem reclamações. Tchau. Ah, não, não pode reclamar, porque senão não tem graça",
		"Eu acho, Elizabeth,  que seria interessante que você olhasse entre os vários cursos que tem sido oferecidos inclusive pelo Senai",
		"Dando um outro exemplo: até agora, a energia hidrelétrica é a mais barata, em termos do que ela dura com a manutenção e também pelo fato da água ser gratuita e da gente poder estocar. O vento podia ser isso também, mas você não conseguiu ainda tecnologia para estocar vento",
		"Então, se a contribuição dos outros países, vamos supor que seja desenvolver uma tecnologia que seja capaz de na eólica estocar, ter uma forma de você estocar, porque o vento ele é diferente em horas do dia. Então, vamos supor que vente mais à noite, como eu faria para estocar isso? Hoje nós usamos as linhas de transmissão, você joga de lá para cá, de lá para lá, para poder capturar isso, mas se tiver uma tecnologia desenvolvida nessa área, todos nós nos beneficiaremos, o mundo inteiro",
		"Ele [o Aedes aegypti] provoca, além da dengue, a chicungunha e ele tem uma variante que transmite o vírus que se chama vírus da zika por causa de uma floresta. Precisamos nos mobilizar para evitar os processos de reprodução do mosquito, porque o mosquito transmite essa doença porque ele coloca o ovo e esse ovo tem o vírus que vai transmitir a doença.",
		"Paes é o prefeito mais feliz do mundo, que dirige a cidade mais importante do mundo e da galáxia. Por que da galáxia? Porque a galáxia é o Rio de Janeiro. A via Láctea é fichinha perto da galáxia que o nosso querido Eduardo Paes tem a honra de ser prefeito.",
		"Quero dizer para vocês que, de fato, Roraima é a capital mais distante de Brasília, mas eu garanto para vocês que essa distância, para nós do Governo Federal, só existe no mapa. E aí eu me considero hoje uma roraimada, roraimada, no que prova que eu estou bem perto de vocês.",
		"Eu acredito que nós teremos uns Jogos Olímpicos que vai ter uma qualidade totalmente diferente e que vai ser capaz de deixar um legado tanto… porque geralmente as pessoas pensam: ‘Ah, o legado é só depois’. Não, vai deixar um legado antes, durante e depois.",
		"Foi muito, houve uma procura imensa, tinham seis empresas que apresentaram suas propostas, houve um deságio de quase… Foi um pouco mais de 38%, mas eu fico em 38% para ninguém dizer: ‘Ah, ela disse que era 38′, mas não é não. É 39, 38 e qualquer coisa ou é 36. 38, eu acho que é 39, mas vou dizer 38.",
		"Não, querido, eu acho que o meu mandato é, eu diria assim, mais firme do que essa rede. Agora, a rede, eu acho que ela tem um lado lúdico, sabe? Porque isso que as crianças gostam tanto no pavilhão. Porque, quando você está lá em cima… Eu não posso ficar aqui brincando, não é? Então… Mas você percebe direitinho como é que dá para brincar, porque se você inclinar para um lado e, imediatamente, virar para o outro, você fica balançando mesmo, você consegue equilibrar.",
		"A 'mosquita', ela, é a 'mosquita' que põe em média 400 ovos. Se você considerar que a 'mosquita' transmite também (o vírus), que é ela que pica, que ela que provoca a contaminação das pessoas.",
		"Nele (no livro), ele diz que nós criamos vínculos sociais e uma das coisas que mais nos une é a fofoca. Uma coisa que nos distingue, que chimpanzé não faz. Orangotango não faz.",
		"O golpe não é só contra mim, é também contra mim, mas não é, sobretudo, contra mim.",
	}
)

func dilma(command *bot.PassiveCmd) (string, error) {
	if re.MatchString(command.Raw) {
		return fmt.Sprintf(":dilma: %s", frasesDilma[rand.Intn(len(frasesDilma))]), nil
	}
	return "", nil
}

func init() {
	bot.RegisterPassiveCommand(
		"dilma",
		dilma)
}
