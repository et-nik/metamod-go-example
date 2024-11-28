package main

import (
	"fmt"
	metamod "github.com/et-nik/metamod-go"
	"github.com/et-nik/metamod-go/engine"
	"github.com/et-nik/metamod-go/vector"
	"math"
	"strconv"
)

var (
	engineFuncs *metamod.EngineFuncs
)

func init() {
	err := metamod.SetPluginInfo(&metamod.PluginInfo{
		InterfaceVersion: metamod.MetaInterfaceVersion,
		Name:             "Metamod Go Example",
		Version:          Version,
		Date:             BuildDate,
		Author:           "KNiK",
		Url:              "https://github.com/et-nik/metamod-go-example",
		LogTag:           "MetamodGoExample",
		Loadable:         metamod.PluginLoadTimeStartup,
		Unloadable:       metamod.PluginLoadTimeAnyTime,
	})
	if err != nil {
		panic(err)
	}

	err = metamod.SetMetaCallbacks(&metamod.MetaCallbacks{
		MetaInit:   MetaInit,
		MetaQuery:  MetaQuery,
		MetaAttach: MetaAttach,
		MetaDetach: func(_ int, _ int) int {
			fmt.Println("Called MetaDetach")

			return 1
		},
	})
	if err != nil {
		panic(err)
	}
}

func main() {}

func MetaInit() {
	fmt.Println()
	fmt.Println("called MetaInit")
	fmt.Println()
}

func MetaQuery() int {
	fmt.Println()
	fmt.Println("called MetaQuery")
	fmt.Println()

	var err error

	engineFuncs, err = metamod.GetEngineFuncs()
	if err != nil {
		fmt.Println("Failed to get engine funcs:", err)
	}

	return 1
}

func MetaAttach(_ int) int {
	fmt.Println()
	fmt.Println("called MetaAttach")
	fmt.Println()

	engineFuncs.AddServerCommand("entinfo", func(argc int, argv ...string) {
		if argc < 2 {
			fmt.Println("Usage: entinfo <entityIndex>")
			return
		}

		entityIndex, err := strconv.Atoi(argv[1])
		if err != nil {
			fmt.Println("Invalid entity index")
			return
		}

		edict := engineFuncs.EntityOfEntIndex(entityIndex)
		if edict == nil {
			fmt.Println("Entity not found")
			return
		}

		if metamod.IsNullEntity(edict) {
			fmt.Println("Entity is null")
			return
		}

		entVars := edict.EntVars()

		if !entVars.IsValid() {
			fmt.Println("Entity is not valid")
			return
		}

		fmt.Println()
		fmt.Println("=====================================")
		fmt.Println("Entity info")
		fmt.Println("Index:", entityIndex)
		fmt.Println("SerialNumber:", edict.SerialNumber())
		fmt.Println("Netname:", edict.EntVars().NetName())
		fmt.Println("Classname:", entVars.ClassName())
		fmt.Println("Globalname:", entVars.GlobalName())
		fmt.Println("Origin:", entVars.Origin())
		fmt.Println("Angles:", entVars.Angles())
		fmt.Println("VAngle:", entVars.VAngle())
		fmt.Println("Model:", entVars.Model())
		fmt.Println("ViewModel:", entVars.ViewModel())
		fmt.Println("WeaponModel:", entVars.WeaponModel())
		fmt.Println("Health:", entVars.Health())
		fmt.Println("Max Health:", entVars.MaxHealth())
		fmt.Println("Max Speed:", entVars.MaxSpeed())
	})

	// Trace line from entity to forward
	engineFuncs.AddServerCommand("traceline", func(argc int, argv ...string) {
		if argc < 3 {
			fmt.Println("Usage: traceline <entityID> <forwardDistance>")

			return
		}

		entityIndex, err := strconv.Atoi(argv[1])
		if err != nil {
			fmt.Println("Invalid entity index")

			return
		}

		edict := engineFuncs.EntityOfEntIndex(entityIndex)
		if edict == nil {
			fmt.Println("Entity not found")

			return
		}

		if metamod.IsNullEntity(edict) {
			fmt.Println("Entity is null")

			return
		}

		entVars := edict.EntVars()

		if !entVars.IsValid() {
			fmt.Println("Entity is not valid")

			return
		}

		distance, err := strconv.ParseFloat(argv[2], 32)
		if err != nil {
			fmt.Println("Invalid distance")

			return
		}

		if distance <= 0 {
			fmt.Println("Distance must be greater than 0")

			return
		}

		if distance > 1000 {
			fmt.Println("Distance must be less than 1000")

			return
		}

		start := entVars.Origin().Add(entVars.ViewOfs())
		end := start.Add(AnglesToForward(entVars.VAngle()).Mul(float32(distance)))

		traceResult := engineFuncs.TraceLine(entVars.Origin(), end, engine.TraceDontIgnoreMonsters, edict)

		fmt.Println()
		fmt.Println("=====================================")
		fmt.Println("Trace line")
		fmt.Println("Start:", start)
		fmt.Println("End:", end)
		fmt.Println("AllSolid:", traceResult.AllSolid)
		fmt.Println("StartSolid:", traceResult.StartSolid)
		fmt.Println("InOpen:", traceResult.InOpen)
		fmt.Println("InWater:", traceResult.InWater)
		fmt.Println("Fraction:", traceResult.Fraction)
		fmt.Println("EndPos:", traceResult.EndPos)

		if !metamod.IsNullEntity(traceResult.Hit) {
			fmt.Println("Hit NetName:", traceResult.Hit.EntVars().NetName())
			fmt.Println("Hit Classname:", traceResult.Hit.EntVars().ClassName())
		} else {
			fmt.Println("Hit: nil")
		}

		fmt.Println("HitGroup:", traceResult.HitGroup)
	})

	return 1
}

// AnglesToForward converts angles to forward vector
func AnglesToForward(angles vector.Vector) vector.Vector {
	pitch := angles[0] * (2 * math.Pi / 360)
	sp := math.Sin(float64(pitch))
	cp := math.Cos(float64(pitch))

	yaw := angles[1] * (2 * math.Pi / 360)
	sy := math.Sin(float64(yaw))
	cy := math.Cos(float64(yaw))

	return vector.Vector{
		float32(cp * cy),
		float32(cp * sy),
		float32(-1 * sp),
	}
}
