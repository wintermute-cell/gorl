package entities

import (
	"gorl/fw/core/entities"
	input "gorl/fw/core/input/input_event"
	"gorl/fw/physics"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Ensure that EnvironmentEntity implements IEntity.
var _ entities.IEntity = &EnvironmentEntity{}

// Environment Entity
type EnvironmentEntity struct {
	*entities.Entity // Required!

	// Custom Fields
	// Add fields here for any state that the entity should keep track of
	// ...
	cols []*physics.Collider // a list of chain shape colliders

	groundTexture rl.Texture2D
}

// NewEnvironmentEntity creates a new instance of the EnvironmentEntity.
func NewEnvironmentEntity() *EnvironmentEntity {
	// NOTE: you can modify the constructor to take any parameters you need to
	// initialize the entity.
	new_ent := &EnvironmentEntity{
		Entity:        entities.NewEntity("EnvironmentEntity", rl.Vector2Zero(), 0, rl.Vector2One()),
		cols:          make([]*physics.Collider, 0),
		groundTexture: rl.LoadTexture("map_thresh_whitewalls_antialiased.png"),
	}

	// a chain collider surrounding the 1920x1080 map
	new_ent.cols = append(new_ent.cols,
		physics.NewChainShapeCollider(
			rl.Vector2Zero(),
			[]rl.Vector2{
				{X: 0, Y: 0}, {X: 1920, Y: 0}, {X: 1920, Y: 1080}, {X: 0, Y: 1080},
			},
			physics.BodyTypeStatic,
		),
	)

	// Add the large colliders for the map geometry
	for _, colPoints := range [][]rl.Vector2{
		collider1Points,
		collider2Points,
		collider3Points,
		collider4Points,
	} {
		new_ent.cols = append(new_ent.cols,
			physics.NewChainShapeCollider(
				rl.Vector2Zero(),
				colPoints,
				physics.BodyTypeStatic,
			),
		)
	}

	// Set the collision category for the colliders
	for _, col := range new_ent.cols {
		col.SetCategory(physics.CollisionCategoryEnvironment)
	}

	return new_ent
}

func (ent *EnvironmentEntity) Init() {
	// Initialization logic for the entity
	// ...
}

func (ent *EnvironmentEntity) Deinit() {
	// De-initialization logic for the entity
	// ...
}

func (ent *EnvironmentEntity) Update() {
	// Update logic for the entity per frame
	// ...
}

func (ent *EnvironmentEntity) Draw() {
	//rl.DrawTexture(ent.groundTexture, 0, 0, rl.White)
	for _, col := range ent.cols {
		verts := col.GetVertices()
		for idx, v := range verts {
			if idx == 0 {
				rl.DrawLineEx(v, verts[len(verts)-1], 2, rl.Green)
			} else {
				rl.DrawLineEx(v, verts[idx-1], 2, rl.Green)
			}
		}
	}
}

func (ent *EnvironmentEntity) OnInputEvent(event *input.InputEvent) bool {
	// Logic to run when an input event is received.
	// Return false if the event was consumed and should not be propagated
	// further.
	// ...
	return true
}

var collider1Points = []rl.Vector2{
	{X: 0, Y: 0}, {X: 320, Y: 0}, {X: 317, Y: 29}, {X: 315, Y: 51},
	{X: 315, Y: 100}, {X: 316, Y: 121}, {X: 317, Y: 150},
	{X: 314, Y: 182}, {X: 310, Y: 203}, {X: 306, Y: 222},
	{X: 302, Y: 244}, {X: 299, Y: 265}, {X: 298, Y: 281},
	{X: 301, Y: 317}, {X: 313, Y: 338}, {X: 331, Y: 351},
	{X: 349, Y: 359}, {X: 370, Y: 367}, {X: 393, Y: 379},
	{X: 406, Y: 398}, {X: 399, Y: 429}, {X: 383, Y: 449},
	{X: 367, Y: 469}, {X: 355, Y: 490}, {X: 347, Y: 507},
	{X: 339, Y: 526}, {X: 331, Y: 543}, {X: 319, Y: 562},
	{X: 303, Y: 581}, {X: 287, Y: 599}, {X: 275, Y: 616},
	{X: 267, Y: 631}, {X: 259, Y: 653}, {X: 257, Y: 684},
	{X: 261, Y: 713}, {X: 273, Y: 735}, {X: 290, Y: 752},
	{X: 311, Y: 768}, {X: 332, Y: 784}, {X: 351, Y: 800},
	{X: 370, Y: 816}, {X: 386, Y: 828}, {X: 401, Y: 836},
	{X: 420, Y: 844}, {X: 447, Y: 849}, {X: 469, Y: 852},
	{X: 508, Y: 850}, {X: 537, Y: 844}, {X: 562, Y: 836},
	{X: 583, Y: 822}, {X: 588, Y: 782}, {X: 585, Y: 756},
	{X: 589, Y: 724}, {X: 601, Y: 700}, {X: 614, Y: 687},
	{X: 628, Y: 679}, {X: 645, Y: 673}, {X: 662, Y: 671},
	{X: 689, Y: 671}, {X: 717, Y: 674}, {X: 738, Y: 677},
	{X: 779, Y: 673}, {X: 800, Y: 665}, {X: 820, Y: 653},
	{X: 842, Y: 641}, {X: 859, Y: 633}, {X: 876, Y: 627},
	{X: 893, Y: 623}, {X: 923, Y: 619}, {X: 947, Y: 607},
	{X: 946, Y: 569}, {X: 930, Y: 547}, {X: 914, Y: 527},
	{X: 898, Y: 505}, {X: 882, Y: 482}, {X: 866, Y: 457},
	{X: 850, Y: 433}, {X: 834, Y: 413}, {X: 818, Y: 393},
	{X: 806, Y: 376}, {X: 798, Y: 360}, {X: 791, Y: 337},
	{X: 793, Y: 293}, {X: 805, Y: 267}, {X: 821, Y: 243},
	{X: 837, Y: 223}, {X: 853, Y: 206}, {X: 871, Y: 190},
	{X: 892, Y: 178}, {X: 912, Y: 170}, {X: 934, Y: 165},
	{X: 969, Y: 168}, {X: 996, Y: 180}, {X: 1019, Y: 196},
	{X: 1037, Y: 212}, {X: 1055, Y: 228}, {X: 1073, Y: 244},
	{X: 1090, Y: 256}, {X: 1105, Y: 264}, {X: 1121, Y: 270},
	{X: 1138, Y: 274}, {X: 1155, Y: 277}, {X: 1171, Y: 278},
	{X: 1204, Y: 282}, {X: 1224, Y: 294}, {X: 1237, Y: 308},
	{X: 1245, Y: 324}, {X: 1250, Y: 341}, {X: 1246, Y: 380},
	{X: 1238, Y: 402}, {X: 1226, Y: 422}, {X: 1207, Y: 439},
	{X: 1185, Y: 451}, {X: 1162, Y: 463}, {X: 1138, Y: 479},
	{X: 1118, Y: 495}, {X: 1101, Y: 511}, {X: 1085, Y: 528},
	{X: 1069, Y: 550}, {X: 1057, Y: 571}, {X: 1049, Y: 590},
	{X: 1047, Y: 612}, {X: 1048, Y: 634}, {X: 1050, Y: 656},
	{X: 1054, Y: 692}, {X: 1052, Y: 729}, {X: 1039, Y: 748},
	{X: 1024, Y: 761}, {X: 1009, Y: 769}, {X: 986, Y: 777},
	{X: 961, Y: 783}, {X: 943, Y: 786}, {X: 924, Y: 787},
	{X: 876, Y: 785}, {X: 850, Y: 779}, {X: 826, Y: 771},
	{X: 804, Y: 763}, {X: 786, Y: 758}, {X: 768, Y: 757},
	{X: 738, Y: 761}, {X: 715, Y: 769}, {X: 698, Y: 777},
	{X: 682, Y: 785}, {X: 667, Y: 793}, {X: 653, Y: 801},
	{X: 631, Y: 813}, {X: 615, Y: 828}, {X: 618, Y: 855},
	{X: 630, Y: 875}, {X: 646, Y: 895}, {X: 662, Y: 913},
	{X: 682, Y: 929}, {X: 703, Y: 941}, {X: 719, Y: 949},
	{X: 740, Y: 957}, {X: 761, Y: 963}, {X: 782, Y: 967},
	{X: 813, Y: 970}, {X: 843, Y: 967}, {X: 863, Y: 963},
	{X: 881, Y: 959}, {X: 898, Y: 955}, {X: 921, Y: 949},
	{X: 949, Y: 943}, {X: 978, Y: 941}, {X: 1010, Y: 938},
	{X: 1032, Y: 932}, {X: 1051, Y: 920}, {X: 1068, Y: 903},
	{X: 1084, Y: 886}, {X: 1100, Y: 869}, {X: 1117, Y: 853},
	{X: 1135, Y: 841}, {X: 1152, Y: 833}, {X: 1171, Y: 828},
	{X: 1213, Y: 829}, {X: 1234, Y: 835}, {X: 1256, Y: 847},
	{X: 1275, Y: 863}, {X: 1287, Y: 884}, {X: 1292, Y: 910},
	{X: 1288, Y: 938}, {X: 1276, Y: 962}, {X: 1259, Y: 980},
	{X: 1236, Y: 996}, {X: 1214, Y: 1008}, {X: 1197, Y: 1016},
	{X: 1181, Y: 1024}, {X: 1161, Y: 1036}, {X: 1142, Y: 1080},
	{X: 1128, Y: 1080}, {X: 197, Y: 1080}, {X: 207, Y: 1080},
	{X: 205, Y: 1011}, {X: 193, Y: 987}, {X: 177, Y: 963},
	{X: 161, Y: 944}, {X: 145, Y: 927}, {X: 128, Y: 911},
	{X: 106, Y: 895}, {X: 83, Y: 883}, {X: 66, Y: 875},
	{X: 47, Y: 863}, {X: 31, Y: 842}, {X: 0, Y: 817},
	{X: 0, Y: 796}, {X: 0, Y: 781}, {X: 0, Y: 325},
	{X: 0, Y: 321}, {X: 53, Y: 315}, {X: 77, Y: 307},
	{X: 102, Y: 295}, {X: 122, Y: 277}, {X: 134, Y: 254},
	{X: 140, Y: 231}, {X: 139, Y: 197}, {X: 133, Y: 171},
	{X: 125, Y: 151}, {X: 117, Y: 135}, {X: 105, Y: 114},
	{X: 89, Y: 92}, {X: 73, Y: 72}, {X: 57, Y: 55},
	{X: 41, Y: 39}, {X: 25, Y: 22},
}

var collider2Points = []rl.Vector2{
	{X: 804, Y: 0}, {X: 1020, Y: 0}, {X: 1006, Y: 20}, {X: 988, Y: 28},
	{X: 969, Y: 33}, {X: 923, Y: 31}, {X: 897, Y: 25}, {X: 873, Y: 17},
	{X: 846, Y: 9}, {X: 820, Y: 3},
}

var collider3Points = []rl.Vector2{
	{X: 1592, Y: 0}, {X: 1861, Y: 0}, {X: 1869, Y: 0}, {X: 1877, Y: 58},
	{X: 1889, Y: 80}, {X: 1920, Y: 92}, {X: 1920, Y: 418}, {X: 1896, Y: 404},
	{X: 1879, Y: 396}, {X: 1842, Y: 397}, {X: 1824, Y: 409}, {X: 1812, Y: 429},
	{X: 1811, Y: 467}, {X: 1815, Y: 492}, {X: 1821, Y: 515}, {X: 1833, Y: 535},
	{X: 1855, Y: 547}, {X: 1891, Y: 540}, {X: 1920, Y: 520}, {X: 1920, Y: 500},
	{X: 1920, Y: 473}, {X: 1920, Y: 826}, {X: 1890, Y: 834}, {X: 1874, Y: 842},
	{X: 1852, Y: 854}, {X: 1828, Y: 870}, {X: 1806, Y: 886}, {X: 1780, Y: 901},
	{X: 1737, Y: 900}, {X: 1718, Y: 892}, {X: 1703, Y: 884}, {X: 1683, Y: 872},
	{X: 1664, Y: 860}, {X: 1650, Y: 852}, {X: 1635, Y: 844}, {X: 1613, Y: 836},
	{X: 1586, Y: 830}, {X: 1559, Y: 824}, {X: 1533, Y: 812}, {X: 1515, Y: 794},
	{X: 1499, Y: 773}, {X: 1483, Y: 752}, {X: 1467, Y: 733}, {X: 1451, Y: 715},
	{X: 1435, Y: 695}, {X: 1423, Y: 676}, {X: 1415, Y: 659}, {X: 1408, Y: 635},
	{X: 1410, Y: 599}, {X: 1422, Y: 578}, {X: 1441, Y: 562}, {X: 1466, Y: 546},
	{X: 1488, Y: 530}, {X: 1507, Y: 513}, {X: 1523, Y: 494}, {X: 1539, Y: 471},
	{X: 1555, Y: 444}, {X: 1567, Y: 422}, {X: 1575, Y: 401}, {X: 1583, Y: 373},
	{X: 1591, Y: 350}, {X: 1603, Y: 328}, {X: 1619, Y: 308}, {X: 1635, Y: 292},
	{X: 1652, Y: 280}, {X: 1666, Y: 272}, {X: 1683, Y: 264}, {X: 1704, Y: 258},
	{X: 1745, Y: 259}, {X: 1764, Y: 262}, {X: 1789, Y: 263}, {X: 1812, Y: 252},
	{X: 1819, Y: 219}, {X: 1816, Y: 196}, {X: 1810, Y: 175}, {X: 1797, Y: 156},
	{X: 1776, Y: 140}, {X: 1752, Y: 124}, {X: 1733, Y: 108}, {X: 1717, Y: 91},
	{X: 1701, Y: 73}, {X: 1685, Y: 55}, {X: 1668, Y: 39}, {X: 1649, Y: 23},
	{X: 1629, Y: 11}, {X: 1592, Y: 0},
}

var collider4Points = []rl.Vector2{
	{X: 1778, Y: 1080}, {X: 1817, Y: 1080}, {X: 1842, Y: 1080},
	{X: 1748, Y: 1080},
}
