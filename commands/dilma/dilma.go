package dilma

import (
	"fmt"
	"github.com/fabioxgn/go-bot"
	"math/rand"
	"regexp"
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
