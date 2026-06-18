# VRC-Shader-Deobfuscator

A lightweight command-line tool written in Go to expand `#define`-based obfuscation in Unity/VRChat HLSL shaders.

Many commercially available Unity/VRChat shaders are distributed with `#define` obfuscation, where names of variables, functions, or characters are mapped to generated strings via `#define` preprocessor directives. This makes debugging, profiling, optimization, and auditing extremely difficult.

This tool ([deobfuscator.go](deobfuscator.go)) acts similarly to a standard C-preprocessor (like `cpp` or `gcc -E`), parsing the `#define` macro mappings, expanding them across the document, and cleaning up the formatting to make the shader readable.

some day i will make a C# based version of this tool, to include it directly into unity

---

## Usage

### Prerequisites
* Go 1.16 or higher installed on your system.

### Running the Deobfuscator
To deobfuscate a shader file, run:
```bash
go run deobfuscator.go <path_to_shader.shader>
```

The tool will:
1. Parse all `#define` directives mapping obfuscated identifiers.
2. Remove the `#define` lines.
3. Replace all obfuscated names with their defined counterparts.
4. Clean up whitespace and format newlines after semicolons (`;`) and braces (`{`, `}`) for readability.
5. Save the output to `<path_to_shader>_deob.shader` (in the same directory as the original).

---

## Legal Notice & Compliance (Rechtlicher Hinweis)
*Please read this section carefully. If you reside in Germany, German copyright law (UrhG) applies to your usage of this software.*

This tool has been developed and is published with strict adherence to German copyright law (**Urheberrechtsgesetz - UrhG**).

### 1. Dual-Use & Preprocessor Equivalence (Standardwerkzeug)
Technically, this program is a simple **macro expander and text formatter**. It performs text search-and-replace operations identical to the preprocessor stage of any standard compiler (such as GCC, Clang, or the Unity shader compiler itself). Since compiler preprocessors are standard, neutral developer utilities, this tool constitutes a neutral, dual-use utility.

### 2. No Circumvention of Access Control / Copy Protection (§ 95a UrhG)
Under **§ 95a UrhG**, it is prohibited to bypass "effective technical protection measures" (*wirksame technische Schutzmaßnahmen*).
* `#define` obfuscation is **not** an access control or copy protection mechanism.
* The obfuscated code is distributed in plaintext within the shader files, and the GPU/compiler must read, parse, and compile these macros in order for the shader to run.
* Because no encryption, decryption, or digital rights management (DRM) bypass is involved, this tool **does not** circumvent any technical protection measures under § 95a UrhG.

### 3. Rights of the Authorized User (§ 69d & § 69e UrhG)
If you hold a valid license to use the shader, German copyright law provides you with specific rights that cannot be contractually excluded (§ 69g Abs. 2 UrhG):
* **Intended Use & Error Correction (§ 69d Abs. 1 UrhG):** The authorized user is permitted to perform acts (such as translation, adaptation, or reproduction) necessary for the intended use of the computer program, including **error correction** (*Fehlerberichtigung*).
* **Interoperability (§ 69e UrhG):** Analysis and modification of the code are permitted if they are necessary to achieve the **interoperability** of an independently created computer program with other programs, provided the information is not already readily available.
* These rights apply to shaders as computer programs under § 69a UrhG.

### 4. Responsibilities of the User
* **License Compliance:** This tool does **not** grant you any licenses or rights to distribute, resell, or share modified shader code. You must ensure you possess a valid license for any shader you analyze.
* **Distribution Ban:** You are responsible for ensuring that your use of the deobfuscated shader code complies with the original shader's license. **Do not redistribute deobfuscated commercial shaders** unless the original license explicitly permits it.
* **Intended Audience:** This tool is intended for personal debugging, performance optimization, security auditing, and educational analysis of shaders you legally own.
