function rmdir(p)
    if not os.isdir(p) then return end

    for _, f in ipairs(os.matchfiles(path.join(p, "**"))) do
        os.remove(f)
    end

    local dirs = os.matchdirs(path.join(p, "**"))
    table.sort(dirs, function(a, b) return #a > #b end)
    for _, d in ipairs(dirs) do
        os.rmdir(d)
    end

    os.rmdir(p)
end


local vendor = "KrabbaTek"

workspace "KrabbaScript"
    configurations { "Debug", "Release" }
    location "build"

    defines { "KSCRIPT_VENDOR=\"" .. vendor .. "\"" }

    project "KrabbaScript"
        kind "ConsoleApp"
        targetname "kscript"

        location "build/krabbascript"

        language "C"
        targetdir "bin/%{prj.name}/%{cfg.buildcfg}"

        includedirs { "include/" }

        files { "src/**.c" }

        filter "configurations:Debug"
            defines { "DEBUG" }
            symbols "On"

        filter "configurations:Release"
            optimize "On"

    project "KrabbaScriptTests"
        kind "ConsoleApp"
        targetname "ktest"

        location "build/ktest"

        language "C"
        targetdir "bin/%{prj.name}/%{cfg.buildcfg}"

        includedirs { "include/", "3rd_party/" }
        defines { "UNITY_INCLUDE_DOUBLE" }

        files { "src/**.c", "tests/**.c", "3rd_party/**.c" }
        removefiles { "src/main.c" }

        filter "configurations:Debug"
            symbols "On"

newaction {
    trigger     = "clean",
    description = "Remove all binaries and generated files",

    execute = function()
        print("Cleaning...")

        os.remove("**.sln")
        os.remove("**.vcxproj")
        os.remove("**.vcxproj.filters")
        os.remove("**.vcxproj.user")
        os.remove("Makefile")
        os.remove("**.make")

        rmdir("bin")
        rmdir("build")

        print("Done.")
    end
}