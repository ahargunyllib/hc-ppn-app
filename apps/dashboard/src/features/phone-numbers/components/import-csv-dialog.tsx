import { Button } from "@/shared/components/ui/button";
import {
  Dialog,
  DialogDescription,
  DialogHeader,
  DialogPanel,
  DialogPopup,
  DialogTitle,
  DialogTrigger,
} from "@/shared/components/ui/dialog";
import { Input } from "@/shared/components/ui/input";
import { toastManager } from "@/shared/components/ui/toast";
import { parseAPIError } from "@/shared/lib/api-client";
import type { CreateUserRequest } from "@/shared/repositories/user/dto";
import { CreateUserSchema } from "@/shared/repositories/user/dto";
import { useCreateUser } from "@/shared/repositories/user/query";
import { Upload } from "lucide-react";
import Papa from "papaparse";
import { useState } from "react";

type CSVRow = {
  phoneNumber: string;
  name: string;
  jobTitle?: string;
  gender?: string;
  dateOfBirth?: string;
};

type ImportResult = {
  total: number;
  success: number;
  failed: number;
  errors: Array<{ row: number; error: string }>;
};

export default function ImportCSVDialog() {
  const [isOpen, setIsOpen] = useState(false);
  const [isImporting, setIsImporting] = useState(false);
  const [progress, setProgress] = useState<string>("");
  const [result, setResult] = useState<ImportResult | null>(null);

  const { mutateAsync } = useCreateUser();

  const handleFileChange = async (
    event: React.ChangeEvent<HTMLInputElement>
  ) => {
    const file = event.target.files?.[0];
    if (!file) {
      return;
    }

    setIsImporting(true);
    setProgress("Membaca file CSV...");
    setResult(null);

    Papa.parse<CSVRow>(file, {
      header: true,
      skipEmptyLines: true,
      transformHeader: (header: string) => header.trim(),
      complete: async (results) => {
        const importResult: ImportResult = {
          total: results.data.length,
          success: 0,
          failed: 0,
          errors: [],
        };

        for (let i = 0; i < results.data.length; i++) {
          const row = results.data[i];
          const rowNumber = i + 2; // +2 because row 1 is header, and index starts at 0

          setProgress(`Mengimpor baris ${rowNumber} dari ${results.data.length + 1}...`);

          try {
            // Prepare the user data
            const userData: CreateUserRequest = {
              phoneNumber: row.phoneNumber?.trim() ?? "",
              name: row.name?.trim() ?? "",
              jobTitle: row.jobTitle?.trim() ?? "",
              gender: row.gender?.trim() ?? "",
              dateOfBirth: row.dateOfBirth?.trim() ?? "",
            };

            // Validate the data
            const validatedData = CreateUserSchema.parse(userData);

            // Create the user
            await mutateAsync(validatedData);

            importResult.success++;
          } catch (error) {
            importResult.failed++;
            const errorMessage =
              error instanceof Error ? error.message : "Unknown error";
            importResult.errors.push({
              row: rowNumber,
              error: errorMessage,
            });
          }
        }

        setResult(importResult);
        setProgress("");
        setIsImporting(false);

        // Show summary toast
        if (importResult.failed === 0) {
          toastManager.add({
            type: "success",
            title: "Import berhasil",
            description: `${importResult.success} user berhasil diimpor.`,
          });
        } else {
          toastManager.add({
            type: "warning",
            title: "Import selesai dengan error",
            description: `${importResult.success} berhasil, ${importResult.failed} gagal.`,
          });
        }

        // Reset file input
        event.target.value = "";
      },
      error: (error) => {
        setIsImporting(false);
        setProgress("");
        toastManager.add({
          type: "error",
          title: "Gagal membaca file CSV",
          description: error.message,
        });
        event.target.value = "";
      },
    });
  };

  const handleClose = () => {
    if (!isImporting) {
      setIsOpen(false);
      setResult(null);
      setProgress("");
    }
  };

  return (
    <Dialog onOpenChange={handleClose} open={isOpen}>
      <DialogTrigger render={<Button size="sm" variant="outline" />}>
        <Upload className="mr-2 h-4 w-4" />
        Import CSV
      </DialogTrigger>
      <DialogPopup>
        <DialogHeader>
          <DialogTitle>Import Users dari CSV</DialogTitle>
          <DialogDescription>
            Upload file CSV untuk menambahkan multiple users sekaligus.
          </DialogDescription>
        </DialogHeader>
        <DialogPanel>
          <div className="flex flex-col gap-4">
            {/* Format Instructions */}
            <div className="rounded-md border border-border bg-muted p-4">
              <h3 className="mb-2 font-semibold text-sm">Format CSV:</h3>
              <div className="space-y-2 text-sm">
                <p>File CSV harus memiliki kolom berikut (dengan header):</p>
                <ul className="ml-4 list-disc space-y-1">
                  <li>
                    <code className="rounded bg-background px-1 py-0.5">
                      phoneNumber
                    </code>{" "}
                    (wajib) - Format E.164, contoh: +6281234567890
                  </li>
                  <li>
                    <code className="rounded bg-background px-1 py-0.5">
                      name
                    </code>{" "}
                    (wajib) - Nama lengkap
                  </li>
                  <li>
                    <code className="rounded bg-background px-1 py-0.5">
                      jobTitle
                    </code>{" "}
                    (opsional) - Jabatan
                  </li>
                  <li>
                    <code className="rounded bg-background px-1 py-0.5">
                      gender
                    </code>{" "}
                    (opsional) - male atau female
                  </li>
                  <li>
                    <code className="rounded bg-background px-1 py-0.5">
                      dateOfBirth
                    </code>{" "}
                    (opsional) - Format: YYYY-MM-DD
                  </li>
                </ul>
                <div className="mt-3">
                  <p className="font-semibold">Contoh:</p>
                  <pre className="mt-1 overflow-x-auto rounded bg-background p-2 font-mono text-xs">
                    {`phoneNumber,name,jobTitle,gender,dateOfBirth
+6281234567890,John Doe,Software Engineer,male,1990-01-15
+6281234567891,Jane Smith,Product Manager,female,1992-05-20`}
                  </pre>
                </div>
              </div>
            </div>

            {/* File Input */}
            <div>
              <Input
                accept=".csv"
                disabled={isImporting}
                onChange={handleFileChange}
                type="file"
              />
            </div>

            {/* Progress */}
            {isImporting && (
              <div className="rounded-md border border-border bg-muted p-3">
                <p className="text-sm">{progress}</p>
              </div>
            )}

            {/* Results */}
            {result && (
              <div className="space-y-2">
                <div className="rounded-md border border-border bg-muted p-3">
                  <h4 className="mb-2 font-semibold text-sm">
                    Hasil Import:
                  </h4>
                  <div className="space-y-1 text-sm">
                    <p>Total: {result.total} baris</p>
                    <p className="text-green-600">Berhasil: {result.success}</p>
                    <p className="text-red-600">Gagal: {result.failed}</p>
                  </div>
                </div>

                {result.errors.length > 0 && (
                  <div className="rounded-md border border-red-200 bg-red-50 p-3">
                    <h4 className="mb-2 font-semibold text-red-800 text-sm">
                      Error Details:
                    </h4>
                    <div className="max-h-40 space-y-1 overflow-y-auto text-sm">
                      {result.errors.map((error, idx) => (
                        <p className="text-red-700" key={idx}>
                          Baris {error.row}: {error.error}
                        </p>
                      ))}
                    </div>
                  </div>
                )}
              </div>
            )}

            <Button disabled={isImporting} onClick={handleClose}>
              {isImporting ? "Mengimpor..." : "Tutup"}
            </Button>
          </div>
        </DialogPanel>
      </DialogPopup>
    </Dialog>
  );
}
