<script lang="ts" setup>
import { h, ref } from 'vue';
import ReverseShell from "./ReverseShell.vue";

const input = ref('');
const visible = ref(false);
const ShowShellTool = ref(false);


const toggleShowShellTool = () => {
  ShowShellTool.value = !ShowShellTool.value;
};

</script>

<template>
 <ReverseShell></ReverseShell>
</template>

<style>
/* Add any styles here */
</style>
