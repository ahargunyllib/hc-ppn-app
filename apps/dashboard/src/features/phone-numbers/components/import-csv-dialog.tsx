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
import { useImportUsersFromCSV } from "@/shared/repositories/user/query";
import { Upload } from "lucide-react";
import { useState } from "react";

type ImportResult = {
  total: number;
  success: number;
  failed: number;
  errors: Array<{ row: number; error: string }>;
};

export default function ImportCSVDialog() {
  const [isOpen, setIsOpen] = useState(false);
  const [result, setResult] = useState<ImportResult | null>(null);

  const { mutate, isPending } = useImportUsersFromCSV();

  const handleFileChange = async (
    event: React.ChangeEvent<HTMLInputElement>
  ) => {
    const file = event.target.files?.[0];
    if (!file) {
      return;
    }

    setResult(null);

    mutate(file, {
      onSuccess: (data) => {
        const importResult: ImportResult = {
          total: data.payload.total,
          success: data.payload.success,
          failed: data.payload.failed,
          errors: data.payload.errors,
        };

        setResult(importResult);

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
      onError: (error) => {
        toastManager.add({
          type: "error",
          title: "Gagal mengimpor CSV",
          description: parseAPIError(error),
        });
        event.target.value = "";
      },
    });
  };

  const handleClose = () => {
    if (!isPending) {
      setIsOpen(false);
      setResult(null);
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
                disabled={isPending}
                onChange={handleFileChange}
                type="file"
              />
            </div>

            {/* Progress */}
            {isPending && (
              <div className="rounded-md border border-border bg-muted p-3">
                <p className="text-sm">Mengimpor data...</p>
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

            <Button disabled={isPending} onClick={handleClose}>
              {isPending ? "Mengimpor..." : "Tutup"}
            </Button>
          </div>
        </DialogPanel>
      </DialogPopup>
    </Dialog>
  );
}
